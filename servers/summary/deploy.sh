sh build.sh

docker push lyons124/summary:latest

ssh ec2-user@api.lyonsupernova.me < update.sh
