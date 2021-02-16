docker rm -f redis
docker run -d --name redis --network info441 redis

docker rm -f gateway
docker pull lyons124/gateway:latest

export TLSCERT=/etc/letsencrypt/live/api.lyonsupernova.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.lyonsupernova.me/privkey.pem
export SESSIONKEY="key"
export MYSQL_ROOT_PASSWORD="password"
export MYSQL_DATABASE="mysqldatabase"
export DSN="root:password@tcp(database:3306)/mysqldatabase"
export REDISADDR=redis:6379

docker run -d \
--name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=:443 \
-e SESSIONKEY=$SESSIONKEY -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=$MYSQL_DATABASE \
-e DSN=$DSN -e REDISADDR=$REDISADDR \
--network info441
lyons124/gateway:latest

exit