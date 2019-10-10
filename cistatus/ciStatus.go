package main

import (
	"encoding/json"
	"github.com/jszwedko/go-circleci"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type CircleCi struct {
	mu       sync.Mutex
	projects map[string]project
	client   *circleci.Client
	gitHubAPI GitHubAPI
}

type project struct {
	reponame      string
	branch        string
	status        string
	url           string
	buildNum      int
	vcsRevision   string
	masterStatus  string
	masterVersion string
}

type Config struct {
	Citoken string
	GHToken string
}


func (circleCi *CircleCi) updateCircleCI(update chan int) {
	for {

		circleCi.updateOnce()

		<-update
	}
}

func (circleCi *CircleCi) readConf() {
	circleCi.client = &circleci.Client{Token: readConfig().Citoken}
	circleCi.gitHubAPI = NewGitHubApi(readConfig().GHToken)
}

func (circleCi *CircleCi) updateOnce() {

	p, e := circleCi.client.ListProjects()
	if e != nil {
		log.Fatal(e)
	}
	circleCi.mu.Lock()
	circleCi.projects = make(map[string]project)
	circleCi.mu.Unlock()
	group := sync.WaitGroup{}
	for _, repo := range p {
		group.Add(1)

		go func(repo *circleci.Project) {
			for _, build := range getBuilds(circleCi.gitHubAPI.context, circleCi.gitHubAPI.client, repo.Reponame) {
				status := "failed"
				if build.TagName != "" {
					status = "success"
				}
				circleCi.mu.Lock()
				circleCi.projects[repo.Reponame+build.Branch] = project{
					reponame:      repo.Reponame,
					branch:        build.Branch,
					status:        status,
					vcsRevision:   build.Commit,
					url: build.ReleaseURL,
				}
				circleCi.mu.Unlock()
			}
			group.Done()
		}(repo)
	}
	group.Wait()
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
