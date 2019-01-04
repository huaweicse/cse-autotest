#!/bin/bash
set -e 
set -x 

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

cd $ROOT_PATH
IMAGE="sdkat_consumer_gosdk"
TAG="latest"

docker build -t $IMAGE:$TAG .
#docker save -o $IMAGE-$TAG.tar $IMAGE:$TAG
