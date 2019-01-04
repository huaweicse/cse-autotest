#!/bin/bash
set -e
set -x

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

pushImageToHuaweiCloudScript="$CURRENT_DIR/push_single_image_to_huaweicloud.sh"

bash "$pushImageToHuaweiCloudScript" "sdkat_consumer_mesher"
bash "$pushImageToHuaweiCloudScript" "sdkat_provider_mesher"