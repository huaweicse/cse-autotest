#!/bin/bash
set -e
set -x

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

docker-compose -f $ROOT_PATH/deployment/docker-compose/docker-compose.yaml up
