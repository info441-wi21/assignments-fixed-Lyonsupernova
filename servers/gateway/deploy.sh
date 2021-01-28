sh build.sh

docker push lyons124/gateway:latest

ssh ec2-user@api.lyonsupernova.me < update.sh
