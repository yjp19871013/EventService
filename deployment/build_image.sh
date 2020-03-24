#!/bin/bash
VERSION=latest
NAME=event_service

IMAGE_NAME=fspub/${NAME}:${VERSION}
CONTAINER_NAME=${NAME}

docker rm -vf ${CONTAINER_NAME}
docker rmi ${IMAGE_NAME}

docker build -t ${IMAGE_NAME} .
docker login
docker push ${IMAGE_NAME}

docker rm -vf ${CONTAINER_NAME}
docker rmi ${IMAGE_NAME}
