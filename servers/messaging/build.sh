GOOS=linux go build

docker build -t lyons124/messaging .

docker push lyons124/messaging

go clean