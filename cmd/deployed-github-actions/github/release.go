package github

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v61/github"
)

// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-a-release-by-tag-name
func GetReleaseByTag(ctx context.Context, client *github.Client, owner string, repo string, tag string) (*github.RepositoryRelease, *github.Response, error) {
	return client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
}

// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
func GetReleaseByLatestTag(ctx context.Context, client *github.Client, owner string, repo string) (*github.RepositoryRelease, *github.Response, error) {
	return client.Repositories.GetLatestRelease(ctx, owner, repo)
}

func DownloadReleaseAsset(ctx context.Context, client *github.Client, owner, repo string, id int64, httpClient *http.Client) (rc io.ReadCloser, redirectURL string, err error) {
	return client.Repositories.DownloadReleaseAsset(ctx, owner, repo, id, httpClient)
}
