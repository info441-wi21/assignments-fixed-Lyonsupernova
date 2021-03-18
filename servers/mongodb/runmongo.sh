docker rm -f info441MongoDB

docker run -d \
    -p 27017:27017 \
    --name info441MongoDB \
    --network info441 \
    mongo