package main

import (
	"fmt"
	"os"
	"time"

	"github.com/deployix/deployed-github-actions/internal/deployed-github-actions/constants"
	promotionsV1 "github.com/deployix/deployed/pkg/promotions/v1"
	utilsV1 "github.com/deployix/deployed/pkg/utils/v1"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type WorkflowInput struct {
	PromotionName string
	Workspace     string
	GitHubToken   string
	SignatureName string
}

func main() {
	// ctx := context.Background()

	input := WorkflowInput{
		PromotionName: os.Getenv("INPUT_PROMOTIONNAME"),
		Workspace:     os.Getenv("GITHUB_WORKSPACE"),
		GitHubToken:   os.Getenv("GITHUB_PAT"),
		SignatureName: constants.GITHUB_SIGNATURE_NAME,
	}

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
	_, err = w.Add(utilsV1.FilePaths().GetChannelsFilePath())
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	commit, err := w.Commit(fmt.Sprintf("Deployed: promote %s", targetedPromotion.Name), &git.CommitOptions{
		Author: &object.Signature{
			Name: input.SignatureName,
			When: time.Now(),
		},
	}) //TODO: custom commit message
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	_, err = r.CommitObject(commit)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

	pushOptions := git.PushOptions{
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: input.SignatureName,
			Password: input.GitHubToken,
		},
	}

	err = r.Push(&pushOptions)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return
	}

}
