#!/bin/bash

docker build -f  Dockerfile.website . -t website
docker build -f  Dockerfile.websocket . -t socket
# docker run -p 80:80 -p 8080:8080 websocket-example:latest 