#!/bin/bash
set -e
set -x

CURRENT_DIR=$(cd $(dirname $0);pwd)
ROOT_PATH=$(dirname $CURRENT_DIR)
export GOPATH="$ROOT_PATH"

cd $ROOT_PATH/src/code.huawei.com/cse/test-case
go test ./...
