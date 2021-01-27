sh build.sh

docker push lyons124/gateway:latest

ssh ec2-user@lyonsupernova.me < update.sh
