package main

import (
	"encoding/json"
	"github.com/jszwedko/go-circleci"
	"github.com/manifoldco/promptui"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	CheckArgs("<environment>")

	r, err := git.PlainOpen(".")
	CheckIfError(err)
	head, err := r.Head()
	CheckIfError(err)

	config, err := r.Config()
	CheckIfError(err)

	url := config.Remotes["origin"].URLs[0]
	index := strings.LastIndex(url, "/")

	repoName := url[index+1 : len(url)-4]

	if os.Args[1] == "prod" {
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
	m["VERSION"] = head.Hash().String()
	m["CIRCLE_JOB"] = "deploy_miljo"
	m["MILJO"] = os.Args[1]

	citoken := readConfig().Citoken
	if len(citoken) == 0 {
		citoken = promtForToken(citoken)
	}
	client := &circleci.Client{Token: citoken}

	build, err := client.ParameterizedBuild("navikt", repoName, head.Name().Short(), m)
	CheckIfError(err)
	Info("Check build status:" + build.BuildURL)
}

func promtForToken(citoken string) string {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "CI token",
		Validate: validate,
	}
	result, err := prompt.Run()
	CheckIfError(err)
	var config = Config{}
	config.Citoken = result
	citoken = result
	confb, e := json.Marshal(config)
	if e != nil {
		log.Fatal(e)
	}
	homeDir, e := os.UserHomeDir()

	CheckIfError(e)
	e = ioutil.WriteFile(homeDir+"/.cistatus.json", confb, 0666)
	CheckIfError(e)
	return citoken
}

func readConfig() Config {
	homeDir, e := os.UserHomeDir()
	if e != nil {
		log.Fatal(e)
	}
	var config = Config{}
	bytes, e := ioutil.ReadFile(homeDir + "/.cistatus.json")
	if e != nil {
		confb, e := json.Marshal(config)
		if e != nil {
			log.Fatal(e)
		}
		e = ioutil.WriteFile(homeDir+"/.cistatus.json", confb, 0666)
		if e != nil {
			log.Fatal(e)
		}
		log.Println("Add Citoken in " + homeDir + "/.cistatus.json")
		os.Exit(-1)
	}

	e = json.Unmarshal(bytes, &config)
	if e != nil {
		log.Println("Add Citoken in " + homeDir + "/.cistatus.json")
		os.Exit(-1)
	}
	return config
}

type Config struct {
	Citoken string
}
