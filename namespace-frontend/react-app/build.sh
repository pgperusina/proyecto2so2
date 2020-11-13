#!/bin/bash
docker login
docker build \
    --no-cache \
    -t react-web-app \
    .
docker tag react-web-app $1/react-web-app
docker push $1/react-web-app
