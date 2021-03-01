sh build.sh

docker push lyons124/messaging

ssh ec2-user@api.lyonsupernova.me < update.sh