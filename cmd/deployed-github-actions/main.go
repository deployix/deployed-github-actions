package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/deployix/deployed-github-actions/cmd/deployed-github-actions/github"
)

func main() {
	ctx := context.Background()

	// STEPS:
	// 1. Get deployed-cli with specific version or default to latest
	input := github.DownloadGithubPackageInput{
		HostURL:   "github.com",                                        //os.Getenv("INPUT_HOST"),
		OrgName:   "deployix",                                          //os.Getenv("INPUT_ORG"),
		RepoName:  "deployed-github-actions",                           //os.Getenv("INPUT_REPO"),
		Version:   "v0.0.1",                                            //os.Getenv("INPUT_VERSION"),
		AssetName: "deployed-github-actions_0.0.1_darwin_amd64.tar.gz", //os.Getenv("INPUT_ASSETNAME"),
	}
	if _, err := github.DownloadPackage(ctx, input); err != nil {
		return
	}
	// 2. Execute deployed-cli with the given arguments
	deployedCLIPath := ""
	cmd := exec.Command(deployedCLIPath, "-h")

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("er")
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))

	// 3. Return outputs

}
