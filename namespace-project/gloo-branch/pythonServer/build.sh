#!/bin/bash
docker login
docker build \
    --no-cache \
    -t python-app \
    .
docker tag python-app $1/python-app
docker push $1/python-app
