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
}

func NewGitHubApi(token string) GitHubAPI {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	location, _ := time.LoadLocation("Europe/Oslo")

	return GitHubAPI{context: ctx, client: client, location: location}
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
	CreatedAt                       time.Time
}

type CurrentDeployments struct {
	prod, q0, q1 *Deployment
}

func (api *GitHubAPI) GetDeployments(repoName string) CurrentDeployments {
	var currentDeployments CurrentDeployments

	for currentPage := 1; ; currentPage++ {
		if lastPage := api.getDeployments(repoName, &currentDeployments, currentPage); lastPage {
			break
		}
	}
	return currentDeployments
}

func (api *GitHubAPI) getDeployments(repoName string, currentDeployments *CurrentDeployments, page int) bool {
	listOptions := github.ListOptions{Page: page}
	deploymentsListOptions := &github.DeploymentsListOptions{ListOptions: listOptions}
	deployments, response, _ := api.client.Repositories.ListDeployments(api.context, owner, repoName, deploymentsListOptions)

	for _, deployment := range deployments {
		var deploymentPayload DeploymentPayload
		if err := json.Unmarshal(deployment.Payload, &deploymentPayload); err != nil {
			continue
		}
		if resources := deploymentPayload.Kubernetes.Resources; len(resources) > 0 {
			currentDeployment := Deployment{
				Environment: *deployment.Environment,
				Namespace:   resources[0].Metadata.Namespace,
				CreatedAt:   deployment.CreatedAt.In(api.location),
				Version:     strings.Split(resources[0].Spec.Image, ":")[1],
			}
			if currentDeployment.Environment == "prod-sbs" && currentDeployments.prod == nil {
				currentDeployments.prod = &currentDeployment
			} else if currentDeployment.Namespace == "q0" && currentDeployments.q0 == nil {
				currentDeployments.q0 = &currentDeployment
			} else if currentDeployment.Namespace == "q1" && currentDeployments.q1 == nil {
				currentDeployments.q1 = &currentDeployment
			}
		}
		if currentDeployments.prod != nil && currentDeployments.q0 != nil && currentDeployments.q1 != nil {
			return true
		}
	}
	return response.LastPage == 0
}
