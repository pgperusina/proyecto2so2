#!/bin/bash
docker login
docker build \
    --no-cache \
    -t go-webserver-grpc \
    .
docker tag go-webserver-grpc $1/go-webserver-grpc
docker push $1/go-webserver-grpc
