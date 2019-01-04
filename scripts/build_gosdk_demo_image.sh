#!/bin/bash
set -e
set -x

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)
export GOPATH="$ROOT_PATH"

ASSETS_DIR="src/code.huawei.com/cse/assets"
BUILD_DIR="build"
BUILD_SCRIPT="build.sh"
BUILD_IMAGE_SCRIPT="build_image.sh"

bash "$ROOT_PATH/$ASSETS_DIR/consumer-gosdk/$BUILD_DIR/$BUILD_SCRIPT"
bash "$ROOT_PATH/$ASSETS_DIR/consumer-gosdk/$BUILD_DIR/$BUILD_IMAGE_SCRIPT"

bash "$ROOT_PATH/$ASSETS_DIR/provider-gosdk/$BUILD_DIR/$BUILD_SCRIPT"
bash "$ROOT_PATH/$ASSETS_DIR/provider-gosdk/$BUILD_DIR/$BUILD_IMAGE_SCRIPT"
