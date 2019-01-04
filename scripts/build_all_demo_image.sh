#!/bin/bash
set -e
set -x

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

bash "$CURRENT_DIR/build_gosdk_demo_image.sh"
bash "$CURRENT_DIR/build_mesher_demo_image.sh"
