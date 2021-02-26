docker rm -f info441MongoDB

docker run -d \
    -p 27017:27017 \
    --name info441MongoDB \
    mongo


docker run --name info441MySQL \
 -e MYSQL_ROOT_PASSWORD="password" \
 -d mysql