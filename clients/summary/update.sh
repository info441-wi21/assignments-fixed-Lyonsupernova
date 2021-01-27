docker rmi -f lyons124/summary:latest
docker pull lyons124/summary:latest
docker rm -f 344summary




docker run -d \
--name 344summary \
-p 80:80 -p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
lyons124/summary:latest

exit