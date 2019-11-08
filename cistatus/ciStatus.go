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

		circleCi.updateOnce()

		<-update
	}
}
func (circleCi *CircleCi) readConf() {
	circleCi.client = &circleci.Client{Token: readConfig().Citoken}
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
			status, _ := circleCi.client.ListRecentBuildsForProject("navikt", repo.Reponame, "", "", 1, 0)
			statusMaster, _ := circleCi.client.ListRecentBuildsForProject("navikt", repo.Reponame, "master", "", 1, 0)
			circleCi.mu.Lock()
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
			circleCi.mu.Unlock()
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
