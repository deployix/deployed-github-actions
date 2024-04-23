package main

import (
	"context"
	"os"

	"github.com/deployix/deployed-github-actions/cmd/deployed-github-actions/github"
)

func main() {
	ctx := context.Background()

	// STEPS:
	// 1. Get deployed-cli with specific version or default to latest
	input := github.DownloadGithubPackageInput{
		HostURL:       os.Getenv("INPUT_HOST"),
		OrgName:       os.Getenv("INPUT_ORG"),
		RepoName:      os.Getenv("INPUT_REPO"),
		Version:       os.Getenv("INPUT_VERSION"),
		FileExtention: os.Getenv("INPUT_"), //TODO: use function to get file extention
	}
	if err := github.DownloadPackage(ctx, input); err != nil {
		return
	}
	// 2. Execute deployed-cli with the given arguments
	// deployedCLIPath := ""
	// cmd := exec.Command(deployedCLIPath, arg0, arg1, arg2, arg3)

	// stdout, err := cmd.Output()
	// if err != nil {
	// 	return
	// }

	// 3. Return outputs

}
