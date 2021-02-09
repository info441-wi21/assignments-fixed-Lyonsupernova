docker rm -f gateway
docker pull lyons124/gateway:latest

export TLSCERT=/etc/letsencrypt/live/api.lyonsupernova.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.lyonsupernova.me/privkey.pem

docker run -d \
--name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
lyons124/gateway:latest

exit