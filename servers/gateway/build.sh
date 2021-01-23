# Builds a Go executable for Linux
GOOS=linux go build
# Builds the Docker container
docker build -t lyons124/gateway .
# Deletes the Go executable for Linux using go clean, so that it doesn't end up getting added to your repo
