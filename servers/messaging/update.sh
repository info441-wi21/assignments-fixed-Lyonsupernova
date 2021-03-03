docker pull lyons124/messaging

docker rm -f messaging

export MONGOADDR="mongodb://mongodb:27017/message"
export PORT=80

# docker rm -f mongodb

# docker run -d \
# -p 27017:27017 \
# --name mongodb \
# --network info441 \
# mongo

# running db instance
docker run -d \
    -e PORT=80 \
    -e MONGOADDR=$MONGOADDR \
    --name messaging \
	--network info441 \
    lyons124/messaging

