docker rm -f redis
docker run -d --name redis --network info441 redis

docker rm -f gateway
docker pull lyons124/gateway:latest

export TLSCERT=/etc/letsencrypt/live/api.lyonsupernova.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.lyonsupernova.me/privkey.pem
export SESSIONKEY=$(openssl rand -base64 18)
export MYSQL_ROOT_PASSWORD="password"
export MYSQL_DATABASE="mysqldatabase"


docker rm -f sqldatabase

docker pull lyons124/sqldatabase

# running db instance
docker run \
    -d \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e MYSQL_DATABASE=$MYSQL_DATABASE \
    --name sqldatabase \
	--network info441 \
    lyons124/sqldatabase

export DSN=root:$MYSQL_DATABASE_PASSWORD@tcp\(sqldatabase:3306\)/$MYSQL_DATABASE
export REDISADDR=redisServer:6379

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
--network info441
lyons124/gateway:latest

exit