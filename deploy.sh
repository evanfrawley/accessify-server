#!/usr/bin/env bash
set -e
source build.sh
docker push evanfrawley/accessify-server

#ssh -oStrictHostKeyChecking=no evan@45.55.211.43 'bash -s' < run.sh
ssh -oStrictHostKeyChecking=no root@api1.accessifyapp.com 'bash -s' < run.sh
