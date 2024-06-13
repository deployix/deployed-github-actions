FROM golang:1.22

ARG DEPLOYED_GITHUB_USERNAME
ARG GITHUB_TOKEN

WORKDIR /deployed-github-actions

COPY . .
RUN echo "machine github.com login ${DEPLOYED_GITHUB_USERNAME} password ${GITHUB_TOKEN}" > ~/.netrc
RUN chmod 600 ~/.netrc

RUN cat ~/.netrc
RUN sleep 10

RUN go env -w GOPRIVATE=github.com/deployix/deployed/*

WORKDIR /deployed-github-actions
RUN go build -o deployed-github-actions cmd/deployed-github-actions/main.go

CMD ["/deployed-github-actions/deployed-github-actions"]