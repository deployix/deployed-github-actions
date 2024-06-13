FROM golang:1.22

WORKDIR /deployed-github-actions

COPY . .


RUN go env -w GOPRIVATE=github.com/deployix/deployed/*

WORKDIR /deployed-github-actions
RUN go build -o deployed-github-actions cmd/deployed-github-actions/main.go

CMD ["/deployed-github-actions/deployed-github-actions"]