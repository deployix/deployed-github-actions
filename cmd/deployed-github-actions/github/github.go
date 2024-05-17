package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/deployix/deployed-github-actions/cmd/deployed-github-actions/constants"
	"github.com/deployix/deployed-github-actions/cmd/deployed-github-actions/utils"
	"github.com/google/go-github/v61/github"
)

type DownloadGithubPackageInput struct {
	HostURL   string
	OrgName   string
	RepoName  string
	AssetName string
	Version   string
}

type DownloadGithubPackageOutput struct {
	File *os.File
}

func DownloadGithubPackage(ctx context.Context, input DownloadGithubPackageInput) (*DownloadGithubPackageOutput, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace("deployix"),
		Password: strings.TrimSpace(""),
	}
	client := github.NewClient(tp.Client())

	// check if the expected verison is latest
	if strings.EqualFold(input.Version, constants.GITHUB_LATEST_RELEASE) {
		//TODO: get latest release and download binary

		release, response, err := GetReleaseByLatestTag(ctx, client, input.OrgName, input.RepoName)
		if err != nil {
			return nil, err
		}
		fmt.Println(release.Name, response.StatusCode)

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
			// download asset locally
			rc, _, err := DownloadReleaseAsset(ctx, client, input.OrgName, input.RepoName, *asset.ID, http.DefaultClient)
			if err != nil {
				return nil, err
			}

			// create local dir/file to store. Use wildcard for prefix so we can keep file extension
			deployedCli, err := os.CreateTemp(constants.DEPLOYED_ARCHIVE_TEMP_PATH, fmt.Sprintf("*-%s", input.AssetName))
			if err != nil {
				return nil, err
			}

			//copy to file
			if _, err = io.Copy(deployedCli, rc); err != nil {
				return nil, err
			}
			defer deployedCli.Close()
			if _, err := deployedCli.Seek(0, 0); err != nil {
				return nil, err
			}

			// extract file
			if err := extractFile(deployedCli, constants.DEPLOYED_CLI_LINUX_PATH); err != nil {
				return nil, err
			}
		}

	}

	return nil, fmt.Errorf("not found")
}

func extractFile(file *os.File, dst string) error {
	fileExtention := filepath.Ext(file.Name())

	switch fileExtention {
	case constants.TAPE_ARCHIVE_FILE_EXTENTION:
		if err := utils.ExtractTarGZ(dst, file); err != nil {
			return err
		}
		break

	case constants.ZIP_ARCHIVE_FILE_EXTENTION:
		break
	default:
		return fmt.Errorf("%s is not a supported archive format", fileExtention)
	}

	return nil
}
