package github

import (
	"context"

	"github.com/google/go-github/v61/github"
)

// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-a-release-by-tag-name
func GetReleaseByTag(ctx context.Context, client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, *github.Response, error) {
	return client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
}
