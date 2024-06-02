FROM golang:1.21

WORKDIR /app
COPY . /app

# Statically compile our app for use in a distroless container
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

COPY . .

ENTRYPOINT ["/app/cmd/deployed-github-actions/main.go"]