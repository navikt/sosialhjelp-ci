package main

import (
	"encoding/json"
	"github.com/jszwedko/go-circleci"
	"io/ioutil"
	"log"
	"os"
)

type CircleCi struct {
	projects  []project
	projectss map[string]project
	client    *circleci.Client
}

type project struct {
	reponame    string
	branch      string
	url         string
	buildNum    int
	vcsRevision string
}

type Config struct {
	Citoken string
}

func (circleCi *CircleCi) update(update chan int) {
	for {
		circleCi.client = &circleci.Client{Token: readConfig().Citoken}
		p, e := circleCi.client.ListProjects()
		if e != nil {
			log.Fatal(e)
		}

		circleCi.projects = make([]project, len(p))
		circleCi.projectss = make(map[string]project)
		for i, repo := range p {
			circleCi.projects[i].reponame = repo.Reponame
			status, _ := circleCi.client.ListRecentBuildsForProject("navikt", repo.Reponame, "", "", 1, 0)
			circleCi.projects[i].branch = status[0].Branch
			circleCi.projects[i].url = status[0].BuildURL
			circleCi.projectss[repo.Reponame] = project{
				reponame:    repo.Reponame,
				branch:      status[0].Branch,
				url:         status[0].BuildURL,
				buildNum:    status[0].BuildNum,
				vcsRevision: status[0].VcsRevision,
			}
		}

		<-update
	}
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
