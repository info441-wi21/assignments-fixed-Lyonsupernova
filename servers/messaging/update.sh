docker pull lyons124/messaging

docker rm -f messaging

export MONGOADDR="mongodb://info441MongoDB:27017/message"
export PORT=80

# running messaging instance
docker run -d \
    -e PORT=80 \
    -e MONGOADDR=$MONGOADDR \
    --name messaging \
	--network info441 \
    lyons124/messaging

