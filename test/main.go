package main

import (
	"encoding/base64"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
	"net/url"
	"time"
)

type ciStatusLayout struct {
	canvas     fyne.CanvasObject
	repolabels []*widget.Hyperlink
	buttons    []*widget.Button
	ciStatus   *CircleCi
	window     fyne.Window
	icon       *fyne.StaticResource

	update chan int
	app    fyne.App
}

func (ci *ciStatusLayout) Layout([]fyne.CanvasObject, fyne.Size) {
	for i, project := range ci.ciStatus.projects {
		u, _ := url.Parse(project.url)
		ci.repolabels[i].Text = fmt.Sprintf("%s\n%s", project.reponame, project.branch)
		ci.repolabels[i].URL = u

		ci.buttons[3*i].Text = "Q0"
		ci.buttons[3*i+1].Text = "Q1"
		ci.buttons[3*i+2].Text = "Prod (master)"
		ci.buttons[3*i].OnTapped = buttonFunc(ci.ciStatus, project.reponame, "Q0")
		ci.buttons[3*i+1].OnTapped = buttonFunc(ci.ciStatus, project.reponame, "Q1")
		ci.buttons[3*i+2].OnTapped = buttonFunc(ci.ciStatus, project.reponame, "Prod (master)")
	}
}

func (ci ciStatusLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(2), objects[0], objects[1])
	return container.Size()
}

func (ci *ciStatusLayout) render() *fyne.Container {

	for _, project := range ci.ciStatus.projects {
		u, _ := url.Parse(project.url)
		hyperlink := widget.NewHyperlink(fmt.Sprintf("%s\n%s", project.reponame, project.branch), u)
		ci.repolabels = append(ci.repolabels, hyperlink)
		ci.buttons = append(ci.buttons, widget.NewButtonWithIcon("", ci.icon, nil))
		ci.buttons = append(ci.buttons, widget.NewButtonWithIcon("", ci.icon, nil))
		ci.buttons = append(ci.buttons, widget.NewButtonWithIcon("", ci.icon, nil))
	}
	var repolabels []fyne.CanvasObject
	for _, value := range ci.repolabels {
		repolabels = append(repolabels, value)
	}
	var buttons []fyne.CanvasObject
	for _, value := range ci.buttons {
		buttons = append(buttons, value)
	}
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		fyne.NewContainerWithLayout(layout.NewGridLayout(1), repolabels...),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3), buttons...))

	container.Layout = ci
	ci.canvas = container

	return container
}

func (ci *ciStatusLayout) animate(canvas fyne.Canvas) {
	tick := time.NewTicker(time.Second)
	tick5s := time.NewTicker(time.Second * 5)
	go func() {
		for {
			ci.Layout(nil, canvas.Size().Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
			canvas.Refresh(ci.canvas)
			<-tick.C
		}
	}()
	go func() {
		for {
			ci.update <- 1
			<-tick5s.C
		}
	}()
}

func buttonFunc(client *CircleCi, reponame, miljo string) func() {
	return func() {
		p := client.projectss[reponame]
		//u, _ := url.Parse(fmt.Sprintf("https://circleci.com/gh/navikt/%s/%d", reponame, p.buildNum+1))
		m := make(map[string]string)
		m["VERSION"] = p.vcsRevision
		m["CIRCLE_JOB"] = "deploy_miljo"
		m["MILJO"] = miljo
		_, e := client.client.ParameterizedBuild("navikt", reponame, p.branch, m)
		if e != nil {
			log.Fatal(e)
		}
		//	e = appl.a.OpenURL(u)
		if e != nil {
			log.Fatal(e)
		}
	}
}

func main() {
	application := app.New()
	ci := &ciStatusLayout{}
	ci.app = application
	ci.window = application.NewWindow("CircleCi deploy app")
	ci.icon = fyne.NewStaticResource("icon", decode(icon))
	application.SetIcon(ci.icon)
	ci.ciStatus = &CircleCi{}
	ci.update = make(chan int)
	go ci.ciStatus.update(ci.update)
	ci.update <- 1
	canvas := ci.render()
	go ci.animate(ci.window.Canvas())
	ci.window.SetContent(canvas)
	ci.window.ShowAndRun()
}

func decode(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
	}
	return data
}

var icon = "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"
