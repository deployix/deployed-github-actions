package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/deployix/deployed-github-actions/cmd/deployed-github-actions/constants"
	"github.com/google/go-github/v61/github"
)

type DownloadGithubPackageInput struct {
	HostURL   string
	OrgName   string
	RepoName  string
	AssetName string
	Version   string
}

func DownloadPackage(ctx context.Context, input DownloadGithubPackageInput) (*github.ReleaseAsset, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace("deployix"),
		Password: strings.TrimSpace(""),
	}
	client := github.NewClient(tp.Client())

	// check if the expected verison is latest
	if strings.EqualFold(input.Version, constants.GITHUB_LATEST_RELEASE) {
		//TODO: get latest release and download binary

	}

	release, response, err := GetReleaseByTag(ctx, client, input.OrgName, input.RepoName, input.Version)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
		return nil, err
	}

	if response.StatusCode == 404 {
		return nil, fmt.Errorf("Unable to find release")
	}

	for _, asset := range release.Assets {
		fmt.Println(*asset.Name)
		if strings.EqualFold(*asset.Name, input.AssetName) {
			fmt.Println("FOUND: " + *asset.Name)
			return asset, nil
		}

	}

	return nil, fmt.Errorf("not found")
}
