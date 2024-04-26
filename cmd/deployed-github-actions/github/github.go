package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v61/github"
)

type DownloadGithubPackageInput struct {
	HostURL       string
	OrgName       string
	RepoName      string
	Version       string
	FileExtention string
}

func DownloadPackage(ctx context.Context, input DownloadGithubPackageInput) (*github.ReleaseAsset, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace("deployix"),
		Password: strings.TrimSpace(""),
	}
	client := github.NewClient(tp.Client())

	release, response, err := GetReleaseByTag(ctx, client, input.OrgName, input.RepoName, input.Version)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
		return nil, err
	}

	for _, asset := range release.Assets {
		fmt.Println(*asset.Name)
		fmt.Println(*asset.BrowserDownloadURL)
	}

	if response.StatusCode == 404 {
		return nil, fmt.Errorf("Unable to find release")
	}

	return nil, nil
}
