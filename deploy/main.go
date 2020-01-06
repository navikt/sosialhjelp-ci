package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jszwedko/go-circleci"
	"github.com/manifoldco/promptui"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	CheckArgs("<environment>\nWhere cwd is a repo and environment is prod | q0 | q1 | dev-gcp | labs-gcp\nThe head ref is matched against tags.")

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
	branch := head.Name().Short()
	version := head.Hash().String()
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
	m := make(map[string]string)
	m["VERSION"] = version
	m["TAG"] = tagName

	conf := readConfig()

	ciClient := getCircleCiClient(conf)
	//	githubClient, ctx := getGithubClient(conf)

	promptConfirm(tagName, environment)

	if environment == "prod" {
		m["CIRCLE_JOB"] = "deploy_prod_tag"
		branch = "master"
	} else if environment == "dev-gcp" || environment == "labs-gcp" { // TODO: Add to help text when ready
		m["CIRCLE_JOB"] = "deploy_gcp"
		m["MILJO"] = environment
	} else {
		m["CIRCLE_JOB"] = "deploy_miljo_tag"
		m["MILJO"] = environment
	}

	build, err := ciClient.ParameterizedBuild("navikt", repoName, branch, m)
	CheckIfError(err)
	Info("Check build status:" + build.BuildURL)
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
