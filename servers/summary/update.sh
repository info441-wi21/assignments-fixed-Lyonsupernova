docker rm -f summary

docker rm -f lyons124/summary

docker pull lyons124/summary

docker run -d \
--name summary \
--network info441 \
lyons124/summary