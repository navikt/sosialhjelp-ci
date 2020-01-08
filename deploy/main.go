package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jszwedko/go-circleci"
	"github.com/manifoldco/promptui"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	CheckArgs("<environment>\nWhere currect working directory is a repo and environment is prod | q0 | q1 | dev-gcp | labs-gcp\nThe head ref is matched against tags.")

	r, err := git.PlainOpen(".")
	CheckIfError(err)

	_ = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})

	head, err := r.Head()
	CheckIfError(err)

	config, err := r.Config()
	CheckIfError(err)

	url := config.Remotes["origin"].URLs[0]
	index := strings.LastIndex(url, "/")
	environment := os.Args[1]
	shouldUseCircleCi := false
	if len(os.Args) > 2 {
		shouldUseCircleCi = os.Args[2] == "circleci"
	}
	branch := head.Name().Short()
	tags, err := r.Tags()
	CheckIfError(err)
	shortHash := head.Hash().String()[:7]
	tagName := ""

	promptForAncestor(branch, r)

	err = tags.ForEach(func(reference *plumbing.Reference) error {
		t := reference.Name().Short()
		if strings.Contains(t, shortHash) {
			tagName = t
		}
		return nil
	})

	CheckIfError(err)
	if len(tagName) == 0 {
		Warning("No tag found, check circleCi")
		os.Exit(1)
	}

	skipNumber := 0
	if strings.Contains(url, ".git") {
		skipNumber = 4
	}
	repoName := url[index+1 : len(url)-skipNumber]

	if environment == "prod" {
		prompt := promptui.Prompt{
			Label:     "Deploy to prod?",
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			os.Exit(0)
			return
		}
		CheckIfError(err)

	}

	promptConfirm(tagName, environment)

	conf := readConfig()
	buildURL := ""
	if shouldUseCircleCi {
		fmt.Println("\nDeployer med CircleCI")
		ciClient := getCircleCiClient(conf)

		m := make(map[string]string)
		m["VERSION"] = head.Hash().String()
		m["TAG"] = tagName
		if environment == "prod" {
			fmt.Println("\nDeployer til PROD")
			m["CIRCLE_JOB"] = "deploy_prod_tag"
			branch = "master"
		} else if environment == "dev-gcp" || environment == "labs-gcp" { // TODO: Add to help text when ready
			fmt.Println("\nDeployer til GCP dev: " + environment)
			m["CIRCLE_JOB"] = "deploy_dev_gcp"
			m["MILJO"] = environment
		} else {
			fmt.Println("\nDeployer til dev: " + environment)
			m["CIRCLE_JOB"] = "deploy_miljo_tag"
			m["MILJO"] = environment
		}

		build, error := ciClient.ParameterizedBuild("navikt", repoName, branch, m)
		buildURL = build.BuildURL
		err = error
	} else {
		fmt.Println("\nDeployer med Github Actions")
		getGithubClient(conf)
		//githubClient, ctx := getGithubClient(conf) // TODO: Bruke githubClient istedenfor dispatch

		dispatch := Dispatch{}
		dispatch.ClientPayload.Tag = tagName

		if environment == "prod" {
			fmt.Println("\nDeployer til PROD")
			dispatch.EventType = "deploy_prod_tag"
		} else if environment == "dev-gcp" || environment == "labs-gcp" { // TODO: Add to help text when ready
			fmt.Println("\nDeployer til GCP dev: " + environment)
			dispatch.EventType = "deploy_dev_gcp"
			dispatch.ClientPayload.Miljo = environment
		} else {
			fmt.Println("\nDeployer til dev: " + environment)
			dispatch.EventType = "deploy_miljo_tag"
			dispatch.ClientPayload.Miljo = environment
		}

		err = createRepositoryDispatch(dispatch, "navikt", repoName)
		buildURL = fmt.Sprintf("https://github.com/%s/%s/actions", "navikt", repoName)
	}
	CheckIfError(err)
	fmt.Println("\nCheck build status: " + buildURL)
}

func promptForAncestor(branch string, r *git.Repository) {
	revHash, err := r.ResolveRevision(plumbing.Revision("origin/" + branch))
	CheckIfError(err)
	revCommit, err := r.CommitObject(*revHash)

	CheckIfError(err)

	headRef, err := r.Head()
	CheckIfError(err)
	headCommit, err := r.CommitObject(headRef.Hash())
	CheckIfError(err)

	isAncestor, err := revCommit.IsAncestor(headCommit)
	CheckIfError(err)
	if !isAncestor {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Head is not updated, are you sure you want to deploy?"),
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			os.Exit(0)
		}
		CheckIfError(err)
	}
}

func promptConfirm(tagName string, environment string) {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Deploy %s to %s?", tagName, environment),
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		os.Exit(0)
	}
	CheckIfError(err)
}

func getCircleCiClient(conf Config) *circleci.Client {
	citoken := conf.Citoken
	if len(citoken) == 0 {
		citoken = promtForCiToken(conf)
	}
	client := &circleci.Client{Token: citoken}
	return client
}

func getGithubClient(conf Config) (*github.Client, context.Context) {
	githubToken := conf.Githubtoken
	if len(githubToken) == 0 {
		githubToken = promtForGithubToken(conf)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}

func promtForCiToken(config Config) string {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "CI token",
		Validate: validate,
	}
	result, err := prompt.Run()
	CheckIfError(err)

	config.Citoken = result
	confb, e := json.Marshal(config)
	if e != nil {
		log.Fatal(e)
	}
	homeDir, e := os.UserHomeDir()

	CheckIfError(e)
	e = ioutil.WriteFile(homeDir+"/.cistatus.json", confb, 0666)
	CheckIfError(e)
	return result
}

func promtForGithubToken(config Config) string {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Github token",
		Validate: validate,
	}
	result, err := prompt.Run()
	CheckIfError(err)

	config.Githubtoken = result
	confb, e := json.Marshal(config)
	if e != nil {
		log.Fatal(e)
	}
	homeDir, e := os.UserHomeDir()

	CheckIfError(e)
	e = ioutil.WriteFile(homeDir+"/.cistatus.json", confb, 0666)
	CheckIfError(e)
	return result
}

func readConfig() Config {
	homeDir, err := os.UserHomeDir()
	CheckIfError(err)

	var config = Config{}
	bytes, err := ioutil.ReadFile(homeDir + "/.cistatus.json")
	if err != nil {
		return config
	}

	err = json.Unmarshal(bytes, &config)
	CheckIfError(err)
	return config
}

type Config struct {
	Citoken     string
	Githubtoken string
}

type Dispatch struct {
	EventType     string        `json:"event_type"`
	ClientPayload ClientPayload `json:"client_payload"`
}

type ClientPayload struct {
	Miljo string `json:"MILJO"`
	Tag   string `json:"TAG"`
}

// TODO: Replace with client library
func createRepositoryDispatch(dispatch Dispatch, owner, repoName string) error {
	dispatchUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/dispatches", owner, repoName)
	payload, e := json.Marshal(dispatch)
	if e != nil {
		return e
	}

	req, _ := http.NewRequest("POST", dispatchUrl, bytes.NewBuffer(payload))
	req.Header.Set("Accept", "application/vnd.github.everest-preview+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", readConfig().Githubtoken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, e := client.Do(req)
	fmt.Println("\nGithub Response: ", resp)
	fmt.Println("\nGithub Error: ", e)
	if e != nil {
		return e
	}
	if e = resp.Body.Close(); e != nil {
		return e
	}
	return nil
}
