package main

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
	"time"
)

const owner = "navikt"

type GitHubAPI struct {
	context  context.Context
	client   *github.Client
	location *time.Location
	deployments map[string]*CurrentDeployments
}

func NewGitHubApi(token string) GitHubAPI {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	location, _ := time.LoadLocation("Europe/Oslo")
	deployments := make(map[string]*CurrentDeployments)

	return GitHubAPI{context: ctx, client: client, location: location, deployments: deployments}
}

type DeploymentPayload struct {
	Kubernetes struct {
		Resources []struct {
			Metadata struct {
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Spec struct {
				Image string `json:"image"`
			} `json:"spec"`
		} `json:"resources"`
	} `json:"kubernetes"`
}

type Deployment struct {
	Environment, Namespace, Version string
	Timestamp                       time.Time
}

type CurrentDeployments struct {
	prod, q0, q1 *Deployment
	Initiated bool
	LastID int64
}

func (api *GitHubAPI) GetDeployments(repoName string, updateDeployments chan bool) {
	for currentPage := 1; ; currentPage++ {
		if lastPage := api.getDeployments(repoName, currentPage); lastPage {
			break
		}
	}
	updateDeployments <- true
}

func (api *GitHubAPI) getDeployments(repoName string, page int) bool {
	deploymentsListOptions := &github.DeploymentsListOptions{ListOptions: github.ListOptions{Page: page}}
	deployments, response, err := api.client.Repositories.ListDeployments(api.context, owner, repoName, deploymentsListOptions)
	if err != nil {
		return true
	}
	currentDeployments, found := api.deployments[repoName]
	if !found {
		currentDeployments = &CurrentDeployments{LastID: 0}
		api.deployments[repoName] = currentDeployments
	}

	for _, deployment := range deployments {
		if currentDeployments.LastID == 0 || *deployment.ID > currentDeployments.LastID {
			currentDeployments.LastID = *deployment.ID
		}
		if currentDeployments.Initiated && *deployment.ID <= currentDeployments.LastID {
			return true
		}

		var deploymentPayload DeploymentPayload
		if err := json.Unmarshal(deployment.Payload, &deploymentPayload); err != nil {
			continue
		}
		if resources := deploymentPayload.Kubernetes.Resources; len(resources) > 0 {
			currentDeployment := Deployment{
				Environment: *deployment.Environment,
				Namespace:   resources[0].Metadata.Namespace,
				Timestamp:   deployment.CreatedAt.In(api.location),
				Version:     strings.Split(resources[0].Spec.Image, ":")[1],
			}
			if currentDeployment.Environment == "prod-sbs" && currentDeployments.prod == nil {
				if success := api.isDeploymentSuccessful(repoName, *deployment.ID); success {
					currentDeployments.prod = &currentDeployment
				}
			} else if currentDeployment.Namespace == "q0" && currentDeployments.q0 == nil {
				if api.isDeploymentSuccessful(repoName, *deployment.ID) {
					currentDeployments.q0 = &currentDeployment
				}
			} else if currentDeployment.Namespace == "q1" && currentDeployments.q1 == nil {
				if api.isDeploymentSuccessful(repoName, *deployment.ID) {
					currentDeployments.q1 = &currentDeployment
				}
			}
		}
		if currentDeployments.prod != nil && currentDeployments.q0 != nil && currentDeployments.q1 != nil {
			currentDeployments.Initiated = true
			return true
		}
	}
	if response.LastPage == 0 {
		currentDeployments.Initiated = true
	}
	return response.LastPage == 0
}

func (api *GitHubAPI) isDeploymentSuccessful(repoName string, deploymentId int64) bool {
	if statuses, _, err := api.client.Repositories.ListDeploymentStatuses(api.context, owner, repoName, deploymentId, nil); err == nil {
		for _, status := range statuses {
			// "queued" -> "in_progress" -> "failure" / ("success" -> "inactive")
			if *status.State == "success" {
				return true
			}
		}
	}
	return false
}
