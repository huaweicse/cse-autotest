#!/bin/bash
set -e
set -x

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

bash "$ROOT_PATH/scripts/build_test_image.sh"
bash "$ROOT_PATH/scripts/build_demo_image.sh"
