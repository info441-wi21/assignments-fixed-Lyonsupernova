GOOS=linux go build

docker build -t lyons124/summary .

docker push lyons124/summary:latest

go clean