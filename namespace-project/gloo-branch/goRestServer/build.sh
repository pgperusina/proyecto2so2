#!/bin/bash
docker login
docker build \
    --no-cache \
    -t go-webserver \
    .
docker tag go-webserver $1/go-webserver
docker push $1/go-webserver
