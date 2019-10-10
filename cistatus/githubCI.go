package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const owner = "navikt"
const digisos = 2442344

type GitHubAPI struct {
	context context.Context
	client *github.Client
}

func NewGitHubApi(token string) GitHubAPI {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return GitHubAPI{context: ctx, client: client}
}

type ReleaseEvent struct {
	Release struct {
		TagName string `json:"tag_name"`
		Commit string `json:"target_commitish"`
	} `json:"release"`
}

type PushEvent struct {
	Ref string       `json:"ref"`
	Commits []struct {
		Commit  string `json:"sha"`
		Message string `json:"message"`
	} `json:"commits"`
}

type Build struct {
	Branch string
	Commit string
	TagName string
	Message string
	ReleaseURL string
}

type ByTypeAndTimestamp []*github.Event

func (repoEvents ByTypeAndTimestamp) Len() int {
	return len(repoEvents)
}

func (repoEvents ByTypeAndTimestamp) Less(i, j int) bool {
	if *repoEvents[i].Type < *repoEvents[j].Type {
		return true
	} else if *repoEvents[i].Type > *repoEvents[j].Type {
		return false
	}
	return repoEvents[i].CreatedAt.After(*repoEvents[j].CreatedAt)
}

func (repoEvents ByTypeAndTimestamp) Swap(i, j int) {
	repoEvents[i], repoEvents[j] = repoEvents[j], repoEvents[i]
}

func getRepoNames(ctx context.Context, client *github.Client) []string {
	teamRepos, _, _ := client.Teams.ListTeamRepos(ctx, digisos, nil)

	var repos []string
	var yesterday = time.Now().Add(-24*time.Hour)

	for _, repo := range teamRepos {
		if repo.UpdatedAt.After(yesterday) {
			repos = append(repos, fmt.Sprintf("%s", *repo.Name))
		}
	}
	return repos
}

func getBuilds(ctx context.Context, client *github.Client, repoName string) []Build {

	repoEvents, _, _ := client.Activity.ListRepositoryEvents(ctx, owner, repoName, nil)

	var builds = map[string]Build{}
	var seenBranch = map[string]bool{}

	sort.Sort(ByTypeAndTimestamp(repoEvents))

	for _, event := range repoEvents {
		if *event.Type == "PushEvent" {
			var pushEvent PushEvent
			err := json.Unmarshal(*event.RawPayload, &pushEvent)

			if err == nil {
				var commits = pushEvent.Commits
				var commit = commits[len(commits)-1].Commit
				var message = commits[len(commits)-1].Message
				var branch = strings.Replace(pushEvent.Ref, "refs/heads/", "", 1)

				if _, seen := seenBranch[branch]; seen {
					continue
				}
				seenBranch[branch] = true

				if _, found := builds[commit]; !found {
					builds[commit] = Build{Commit: commit, Message: message, Branch: branch}
				}
			}
		} else if *event.Type == "ReleaseEvent" {
			var releaseEvent ReleaseEvent
			err := json.Unmarshal(*event.RawPayload, &releaseEvent)

			if err == nil {
				var commit = releaseEvent.Release.Commit
				var tagName = releaseEvent.Release.TagName

				if val, found := builds[commit]; found {
					val.TagName = tagName
					val.ReleaseURL = fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s", owner, repoName, tagName)
					builds[commit] = val
				}
			}
		}
	}

	var buildList []Build

	for _, build := range builds {
		buildList = append(buildList, build)
	}

	return buildList
}

// TODO: Replace with API call once implemented
func createRepositoryDispatch(ci *ciStatusLayout, reponame, miljo string) error {
	dispatchUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/dispatches", owner, reponame)

	var payload= []byte(fmt.Sprintf(`{"event_type": "%s"}`, miljo))
	req, _ := http.NewRequest("POST", dispatchUrl, bytes.NewBuffer(payload))
	req.Header.Set("Accept", "application/vnd.github.everest-preview+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", readConfig().GHToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		return e
	}
	if e = resp.Body.Close(); e != nil {
		return e
	}

	u, _ := url.Parse(fmt.Sprintf("https://github.com/%s/%s/actions", owner, reponame))
	if e = ci.app.OpenURL(u); e != nil {
		return e
	}
	return nil
}
