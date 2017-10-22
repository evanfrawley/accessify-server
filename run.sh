#!/usr/bin/env bash
docker rm -f accessify-server
docker pull evanfrawley/accessify-server

export TLSCERT=/etc/letsencrypt/live/api.accessifyapp.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.accessifyapp.com/privkey.pem

docker run -d \
--name accessify-server \
-v /root/keys:/keys:ro \
-p 80:80 \
evanfrawley/accessify-server

