docker network rm info441
docker network create info441

docker rm -f redis
docker run -d --name redis --network info441 redis

# redeploy summary microservice
docker rm -f summary
docker rm -f lyons124/summary
docker pull lyons124/summary
docker run -d \
    --name summary \
    --network info441 \
    lyons124/summary

# redeploy messaging microservice
docker pull lyons124/messaging
docker rm -f messaging
export MONGOADDR="mongodb://info441MongoDB:27017/message"
export PORT=80
docker run -d \
    -e PORT=80 \
    -e MONGOADDR=$MONGOADDR \
    --name messaging \
	--network info441 \
    lyons124/messaging


docker rm -f gateway
docker pull lyons124/gateway:latest

export TLSCERT=/etc/letsencrypt/live/api.lyonsupernova.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.lyonsupernova.me/privkey.pem
export SESSIONKEY=$(openssl rand -base64 18)
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
export DB_NAME=441sqlserver


docker rm -f 441mysql

docker pull lyons124/441mysql

# running db instance
docker run -d \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e MYSQL_DATABASE=$DB_NAME \
    --name 441mysql \
	--network info441 \
    lyons124/441mysql

export DSN=root:$MYSQL_ROOT_PASSWORD@tcp\(441mysql:3306\)/$DB_NAME
export REDISADDR=redis:6379
export MESSAGESADDR=messaging:80
export SUMMARYADDR=summary:80

docker run -d \
--name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=:443 \
-e SESSIONKEY=$SESSIONKEY \
-e DSN=$DSN \
-e REDISADDR=$REDISADDR \
-e MESSAGESADDR=$MESSAGESADDR \
-e SUMMARYADDR=$SUMMARYADDR \
--network info441 \
lyons124/gateway:latest

exit