package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	constantsV1 "github.com/deployix/deployed/pkg/constants/v1"
	promotionsV1 "github.com/deployix/deployed/pkg/promotions/v1"
	utilsV1 "github.com/deployix/deployed/pkg/utils/v1"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/tcnksm/go-gitconfig"
)

type WorkflowInput struct {
	PromotionName string
	Workspace     string
}

func main() {
	// ctx := context.Background()

	input := WorkflowInput{
		PromotionName: "local-to-dev", //os.Getenv("INPUT_PROMOTIONNAME"),
		Workspace:     os.Getenv("GITHUB_WORKSPACE"),
	}

	// set working directory for filepath
	os.Setenv(constantsV1.FILEPATH_WORKING_DIR_ENV, input.Workspace)

	// get promotions
	promotions, err := promotionsV1.GetPromotions()
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}

	// verify promotion expected exists
	if !promotions.PromotionExists(input.PromotionName) {
		fmt.Printf("promotion with the name `%s` does not exist", input.PromotionName)
		return
	}

	// promote targeted promotion resource
	targetedPromotion := promotions.Promotions[input.PromotionName]
	if err := targetedPromotion.Promote(); err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	// commit changes to default branch
	r, err := git.PlainOpen(input.Workspace)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	w, err := r.Worktree()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	// add channels.yml file to commit as thats what has changed
	filePath := filepath.Join(input.Workspace, utilsV1.FilePaths.GetChannelsFilePath())
	_, err = w.Add(filePath)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	// print git status
	status, err := w.Status()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println(status)

	username, err := gitconfig.Username()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	email, err := gitconfig.Email()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	url, err := gitconfig.OriginURL()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	commit, err := w.Commit(fmt.Sprintf("Deployed: promote %s", targetedPromotion.Name), &git.CommitOptions{
		Author: &object.Signature{
			Name:  username,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	_, err = r.CommitObject(commit)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	repo, err := gitconfig.Repository()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}
	pushOptions := git.PushOptions{
		RemoteName: repo,
		RemoteURL:  url,
	}

	// Set token auth if passed in
	if os.Getenv("INPUT_GITHUBTOKEN") != "" {
		pushOptions.Auth = &http.TokenAuth{
			Token: os.Getenv("INPUT_GITHUBTOKEN"),
		}
	}

	err = r.Push(&pushOptions)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

}
