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
		PromotionName: os.Getenv("INPUT_PROMOTIONNAME"),
		Workspace:     os.Getenv("GITHUB_WORKSPACE"),
	}

	fmt.Println(fmt.Sprintf("promotion name: %s", input.PromotionName))

	// set working directory for filepath
	os.Setenv(constantsV1.FILEPATH_WORKING_DIR_ENV, input.Workspace)
	fmt.Println(fmt.Sprintf("Working Dir: %s", input.Workspace))

	fmt.Println("getting promotions")
	// get promotions
	promotions, err := promotionsV1.GetPromotions()
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}

	fmt.Println("checking promotions exists")

	// verify promotion expected exists
	if !promotions.PromotionExists(input.PromotionName) {
		fmt.Printf("promotion with the name `%s` does not exist", input.PromotionName)
		return
	}

	fmt.Println("updating promotion")
	// promote targeted promotion resource
	targetedPromotion := promotions.Promotions[input.PromotionName]
	if err := targetedPromotion.Promote(); err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("commiting change")
	// commit changes to default branch
	r, err := git.PlainOpen(input.Workspace)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("worktree")
	w, err := r.Worktree()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("add channel.yml " + utilsV1.FilePaths.GetChannelsFilePath())
	// add channels.yml file to commit as thats what has changed
	filePath := filepath.Join(input.Workspace, utilsV1.FilePaths.GetChannelsFilePath())
	_, err = w.Add(filePath)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("status")
	// print git status
	status, err := w.Status()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println(status)

	fmt.Println("username")
	username, err := gitconfig.Username()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("email")
	email, err := gitconfig.Email()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("origin")
	_, err = gitconfig.OriginURL()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("commit object")
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

	fmt.Println("commit")
	_, err = r.CommitObject(commit)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	fmt.Println("repo")
	_, err = gitconfig.Repository()
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}
	pushOptions := git.PushOptions{
		Progress: os.Stdout,
	}

	// Set token auth if passed in
	if os.Getenv("INPUT_GITHUBTOKEN") != "" {
		pushOptions.Auth = &http.TokenAuth{
			Token: os.Getenv("INPUT_GITHUBTOKEN"),
		}
	}

	fmt.Println("push")
	err = r.Push(&pushOptions)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

}
