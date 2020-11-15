#!/bin/bash
docker login
docker build \
    -t go-webserver-grpc2 \
    .
docker tag go-webserver-grpc2 $1/go-webserver-grpc2
docker push $1/go-webserver-grpc2
