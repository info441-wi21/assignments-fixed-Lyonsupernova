docker rm -f database

docker pull lyons124/sqldatabase

docker run \
    -d \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD="password" \
    -e MYSQL_DATABASE="mysqldatabase" \
    --name database \
	--network info441 \
    lyons124/sqldatabase