#!/bin/bash
docker login
docker build \
    --no-cache \
    -t python-app-grpc \
    .
docker tag python-app-grpc $1/python-app-grpc
docker push $1/python-app-grpc
