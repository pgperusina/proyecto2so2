#!/bin/bash
docker login
docker build \
    -t go-webserver-grpc \
    .
docker tag go-webserver-grpc $1/go-webserver-grpc
docker push $1/go-webserver-grpc
