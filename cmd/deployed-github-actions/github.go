package main

import (
	"net/http"

	"github.com/google/go-github/v60/github"
)

type NewGithubClientInput struct {
	HttpClient *http.Client
}

func NewGitHubClient(input NewGithubClientInput) *github.Client {
	return github.NewClient(nil)
}
