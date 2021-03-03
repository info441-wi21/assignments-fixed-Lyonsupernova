docker rm -f message

docker rm -f lyons124/messaging

docker pull lyons124/messaging

export MONGOADDR="mongodb://mongodb:27017/message"
export PORT=80
export MONGOPORT=mongodb:27017


export MYSQL_PORT=3306

# running db instance
docker run -d \
    -e PORT=80 \
    -e MYSQL_PORT=3306 \
    -e MONGOPORT=$MONGOPORT \
    --name message \
	--network info441 \
    lyons124/messaging

