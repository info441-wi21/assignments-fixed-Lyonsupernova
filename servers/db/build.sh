GOOS=linux go build
docker build -t lyons124/sqldatabase .
go clean

docker push lyons124/sqldatabase

ssh ec2-user@api.lyonsupernova.me < deploy.sh 
