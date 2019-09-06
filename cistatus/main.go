package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/jszwedko/go-circleci"
	"github.com/skratchdot/open-golang/open"
	"log"
	"os"
	"time"
)

type App struct {
	projects []*circleci.Project
	a        fyne.App
	w        fyne.Window
}

var appl = App{}

func main() {
	appl.a = app.New()

	appl.w = appl.a.NewWindow("CircleCi deploy app")
	go updateProjects()
	appl.w.ShowAndRun()
}

func updateProjects() {
	for {
		client := &circleci.Client{Token: os.Getenv("CITOKEN")}
		projects, e := client.ListProjects()
		appl.projects = projects
		for e != nil {
			time.Sleep(10 * time.Second)
			appl.projects, e = client.ListProjects()
		}

		var containers []fyne.CanvasObject
		for _, repo := range projects {
			reponame := repo.Reponame
			containers = append(containers, fyne.NewContainer(
				widget.NewLabel(reponame)),
				widget.NewButton("Q0", func() {
					status, _ := client.ListRecentBuildsForProject("navikt", reponame, "", "", -1, 0)

					m := make(map[string]string)
					m["VERSION"] = status[0].VcsRevision
					m["CIRCLE_JOB"] = "deploy_miljo"
					m["MILJO"] = "q0"
					client.ParameterizedBuild("navikt", reponame, status[0].Branch, m)

					e := open.Run(fmt.Sprintf("https://circleci.com/gh/navikt/%s/%d", reponame, status[0].BuildNum+1))
					if e != nil {
						log.Fatal(e)
					}
				}),
				widget.NewButton("Q1", func() {
					status, _ := client.ListRecentBuildsForProject("navikt", reponame, "", "", -1, 0)
					m := make(map[string]string)
					m["VERSION"] = status[0].VcsRevision
					m["CIRCLE_JOB"] = "deploy_miljo"
					m["MILJO"] = "q1"
					client.ParameterizedBuild("navikt", reponame, status[0].Branch, m)

					e := open.Run(fmt.Sprintf("https://circleci.com/gh/navikt/%s/%d", reponame, status[0].BuildNum+1))
					if e != nil {
						log.Fatal(e)
					}
				}),
				widget.NewButton("Prod", func() {
					status, _ := client.ListRecentBuildsForProject("navikt", reponame, "", "", -1, 0)
					m := make(map[string]string)
					m["VERSION"] = status[0].VcsRevision
					m["CIRCLE_JOB"] = "deploy_miljo"
					m["MILJO"] = "deploy_prod"
					client.ParameterizedBuild("navikt", reponame, status[0].Branch, m)

					e := open.Run(fmt.Sprintf("https://circleci.com/gh/navikt/%s/%d", reponame, status[0].BuildNum+1))
					if e != nil {
						log.Fatal(e)
					}
				}),
			)
		}
		containers = append(containers, widget.NewButton("Quit", func() {
			appl.a.Quit()
		}))

		appl.w.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(4), containers...))

		time.Sleep(10 * time.Second)
	}
}
