package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GH_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Cannot get access token from environment. You must set evnvar GH_ACCESS_TOKEN.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	page := 0
	perPage := 100
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	query := "topic:mattermost-plugin"

	repositories := []github.Repository{}
	for {
		ret, resp, err := client.Search.Repositories(ctx, query, opt)
		if err != nil {
			log.Fatalf("Failed to search repositories: %s", err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to request: %s", resp.Status)
		}
		repositories = append(repositories, ret.Repositories...)
		if ret.GetTotal() < perPage {
			break
		}
		opt.Page++
	}
	list := []string{}
	for _, repo := range repositories {
		list = append(list, fmt.Sprintf("github.com/%s", repo.GetFullName()))
	}
	if err := ioutil.WriteFile("topic-repositories.txt", []byte(strings.Join(list, "\n")), os.ModePerm); err != nil {
		log.Fatalf("Failed to write file: %s", err.Error())
	}
}
