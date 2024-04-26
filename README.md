# deployed-github-actions
deployed-github-actions

Example: https://medium.com/@yanzay/building-github-actions-using-go-80a0add54104

# Go release commands
- git tag -a v0.0.1 -m "First release"
- git push origin v0.0.1
- export GITHUB_TOKEN=[TOKEN]
- goreleaser release