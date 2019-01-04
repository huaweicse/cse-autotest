#!/bin/bash
set -e
set -x

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

bash "$CURRENT_DIR/push_gosdk_demo_to_huaweicloud.sh"
bash "$CURRENT_DIR/push_mesher_demo_to_huaweicloud.sh"
