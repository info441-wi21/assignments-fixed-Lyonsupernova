sh build.sh
sh ../messaging/build.sh
sh ../summary/build.sh

ssh ec2-user@api.lyonsupernova.me < update.sh
