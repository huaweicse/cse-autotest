#!/bin/bash
set -e
set -x

if [ ! -z "$SDKAT_SWR_LOGIN_CMD" ]; then
    echo "------------------------Login SWR------------------------------------"
    $SDKAT_SWR_LOGIN_CMD
else
    echo "No SDKAT_SWR_LOGIN_CMD set, will not login in SWR"
fi

echo "--------------------------Push image-----------------------------"
mesherVersion="${SDKAT_MESHER_VERSION:-latest}"
imageLatest="${SDKAT_SWR_ADDR:-100.125.0.198:20202}/${SDKAT_SWR_ORG:-hwcse}/cse-mesher:latest"
docker tag "cse-mesher:${mesherVersion}" "$imageLatest"
docker push "$imageLatest"

imageName="${SDKAT_SWR_ADDR:-100.125.0.198:20202}/${SDKAT_SWR_ORG:-hwcse}/cse-mesher:$mesherVersion"
docker tag "cse-mesher:${mesherVersion}" "$imageName"
docker push "$imageName"


