package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
)

type DownloadGithubPackageInput struct {
	HostURL       string
	OrgName       string
	RepoName      string
	Version       string
	FileExtention string
}

func downloadPackage(ctx context.Context, input DownloadGithubPackageInput) error {
	client := github.NewClient(nil)
	// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-a-release-by-tag-name
	_, response, err := client.Repositories.GetReleaseByTag(ctx, input.OrgName, input.RepoName, input.Version)
	if err != nil {
		return err
	}

	if response.StatusCode == 404 {
		return fmt.Errorf("Unable to find release")
	}

	return nil
}
