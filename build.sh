#!/usr/bin/env bash
GOOS=linux go build
docker build -t evanfrawley/accessify-server .
go clean
