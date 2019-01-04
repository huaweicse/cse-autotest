#!/bin/bash

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)

if [[ -z "$1" || "$(basename $1)" != "auth.yaml" ]]; then
    echo "Please input your auth.yaml path!"
    exit 1
fi
sourceAuthFile=$1

cd $ROOT_PATH

ASSETS_DIR="$ROOT_PATH/src/code.huawei.com/cse/assets"
authFileList=$(find "$ASSETS_DIR" -name "auth.yaml")

for data in ${authFileList[@]}
do
    targetDir=$(dirname $data)
    cp $sourceAuthFile $targetDir
done

echo "Replace auth.yaml success!"
