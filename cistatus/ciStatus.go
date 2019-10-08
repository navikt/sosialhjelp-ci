package main

import (
	"encoding/json"
	"github.com/jszwedko/go-circleci"
	"io/ioutil"
	"log"
	"os"
)

type CircleCi struct {
	projects map[string]project
	client   *circleci.Client
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
}

func (circleCi *CircleCi) update(update chan int) {
	for {
		circleCi.client = &circleci.Client{Token: readConfig().Citoken}
		p, e := circleCi.client.ListProjects()
		if e != nil {
			log.Fatal(e)
		}

		circleCi.projects = make(map[string]project)
		for _, repo := range p {
			status, _ := circleCi.client.ListRecentBuildsForProject("navikt", repo.Reponame, "", "", 1, 0)
			statusMaster, _ := circleCi.client.ListRecentBuildsForProject("navikt", repo.Reponame, "master", "", 1, 0)
			circleCi.projects[repo.Reponame] = project{
				reponame:      repo.Reponame,
				branch:        status[0].Branch,
				status:        status[0].Status,
				url:           status[0].BuildURL,
				buildNum:      status[0].BuildNum,
				vcsRevision:   status[0].VcsRevision,
				masterStatus:  statusMaster[0].Status,
				masterVersion: statusMaster[0].VcsRevision,
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
