#!/bin/bash
set -e
set -x

PKG_NAME="source.tar.gz"

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

cd $ROOT_PATH

if [ -f $PKG_NAME ]; then
    rm $PKG_NAME
fi

tar -zcf $PKG_NAME --exclude src/code.huawei.com/cse/assets src

IMAGE="sdkat_testcase"
TAG="latest"

docker build -t $IMAGE:$TAG .

echo "Build success!"
