#!/bin/bash
set -e
set -x

if [ -z "$1" ];then
    echo "please input image name(with no tag)"
    exit 1
fi

if [ -z "$SDKAT_BUILD_VERSION" ]; then
    echo "No build version set, please set env SDKAT_BUILD_VERSION"
    echo "Use image tag: latest"
fi

if [ ! -z "$SDKAT_SWR_LOGIN_CMD" ]; then
    echo "------------------------Login SWR------------------------------------"
    $SDKAT_SWR_LOGIN_CMD
else
    echo "No SDKAT_SWR_LOGIN_CMD set, will not login in SWR"
fi

echo "--------------------------Push image-----------------------------"
imageLatest="${SDKAT_SWR_ADDR:-100.125.0.198:20202}/${SDKAT_SWR_ORG:-hwcse}/$1:latest"
docker tag "$1:latest" "$imageLatest"
docker push "$imageLatest"

imageName="${SDKAT_SWR_ADDR:-100.125.0.198:20202}/${SDKAT_SWR_ORG:-hwcse}/$1:${SDKAT_BUILD_VERSION:-latest}"
docker tag "$1:latest" "$imageName"
docker push "$imageName"
