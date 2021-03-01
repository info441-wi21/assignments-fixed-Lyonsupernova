docker rm -f message

docker rm -f lyons124/messaging

docker pull lyons124/messaging

export MONGOADDR="mongodb://mongodb:27017/message"
export PORT=80
export MONGOPORT=mongodb:27017

export MYSQL_ROOT_PASSWORD="password"
export MYSQL_DATABASE="mysqldatabase"
export MYSQL_HOST="info441"
export MYSQL_USER="root"
export MYSQL_PORT=3306

# running db instance
docker run \
    -d \
    -p 3306:3306 \
    -e MYSQL_HOST=$MYSQL_HOST \
    -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    -e MySQL_USER=$MYSQL_USER \
    -e MYSQL_DATABASE=$MYSQL_DATABASE \
    -e PORT=80 \
    -e MYSQL_PORT=3306 \
    -e MONGOPORT=$MONGOPORT \
    --name message \
	--network info441 \
    lyons124/messaging

