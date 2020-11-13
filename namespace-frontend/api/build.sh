#!/bin/bash
docker login
docker build \
    --no-cache \
    -t nodejs-api \
    .
docker tag nodejs-api $1/nodejs-api
docker push $1/nodejs-api
