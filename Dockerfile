FROM golang:1.22

ARG GITHUB_LOGIN
ARG GITHUB_TOKEN

WORKDIR /deployed-github-actions

COPY . .
RUN echo "machine github.com login deployed-ci password ${GITHUB_TOKEN}" > ~/.netrc
RUN chmod 600 ~/.netrc

RUN go env -w GOPRIVATE=github.com/deployix/deployed/*


WORKDIR /deployed-github-actions
RUN go build -o deployed-github-actions cmd/deployed-github-actions/main.go

CMD ["/deployed-github-actions/deployed-github-actions"]