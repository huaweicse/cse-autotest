#!/bin/bash

#SDKAT_TENANT_NAME:       tenant name
#SDKAT_USER_NAME:         user name
#SDKAT_PASSWORD:          password
#SDKAT_MESHER_VERSION:    mesher version
#SDKAT_BUILD_VERSION:     build version, it is a time stamp, like 20181130152417
#SDKAT_STACK_ID:          the stack to update
#SDKAT_IAM_ADDR:          iam addr, no scheme
#SDKAT_AOS_ADDR:          aos addr, no scheme
#SDKAT_REGION:            region
#SDKAT_SWR_ADDR:          swr addr
#SDKAT_SWR_ADDR_INTERNAL  some times the outside swr addr has bad certificate, so the program in the cluster can only use internal addr
#SDKAT_SWR_ORG:           swr orgnization
#SDKAT_SWR_LOGIN_CMD:     swr login cmd

set -e
set -x

# tenent info
TenantName="${SDKAT_TENANT_NAME}"
UserName="${SDKAT_USER_NAME}"
PassWd="${SDKAT_PASSWORD}"
if [[ -z "$TenantName" || -z "$UserName" || -z "$PassWd" ]]; then
  echo "Please input valid tenant nane/user name/password"
  exit 1
fi

stackId="${SDKAT_STACK_ID}"
if [ -z "$stackId" ]; then
  echo "Please input valid stack id"
  exit 1
fi

# huawei cloud address
IAMAddr="${SDKAT_IAM_ADDR:-iam.cn-north-1.myhuaweicloud.com}"
AosAddr="${SDKAT_AOS_ADDR:-aos.cn-north-1.myhuaweicloud.com}"
RegionName="${SDKAT_REGION:-cn-north-1}"

# api info
Scheme="https://"
TokenHeaderKey="X-Auth-Token"
StackLifeApi="/v2/stacks/${stackId}/actions"


echo "--------------------------get token from iam---------------------------------"
tmpToken=$(curl -k -s -D - -o /dev/null -X POST \
  https://$IAMAddr/v3/auth/tokens \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d "{ 
  \"auth\": { 
    \"identity\": { 
      \"methods\": [ 
        \"password\" 
      ], 
      \"password\": { 
        \"user\": { 
          \"name\": \"$UserName\", 
          \"password\": \"$PassWd\", 
          \"domain\": { 
            \"name\": \"$TenantName\" 
          } 
        } 
      } 
    }, 
    \"scope\": { 
      \"project\": { 
        \"name\": \"$RegionName\"
      } 
    } 
  } 
}" | grep "X-Subject-Token" | awk -F ": " '{print $2}'|tr "\r" " "|awk '{$1=$1;print}')


echo "--------------------------update stack---------------------------------"

swrAddr="${SDKAT_SWR_ADDR_INTERNAL:-$SDKAT_SWR_ADDR}"
getFullImage() {
    echo "${swrAddr:-100.125.0.198:20202}/${SDKAT_SWR_ORG:-hwcse}/$1:${SDKAT_BUILD_VERSION:-latest}"
}

#stack info
MesherImage="${swrAddr:-100.125.0.198:20202}/${SDKAT_SWR_ORG:-hwcse}/cse-mesher:${SDKAT_MESHER_VERSION:-latest}"

if [ -z "$1" ];then
    echo "Upgrade all the images, if only want to upgrade images related with sdk type, please input gosdk/mesher"
    curl -k -X PUT -H "$TokenHeaderKey:$tmpToken" "$Scheme$AosAddr$StackLifeApi" -d "{
  \"lifecycle\":\"upgrade\",
  \"inputs\":{
    \"consumer_gosdk_image_addr\":\"$(getFullImage sdkat_consumer_gosdk)\",
    \"provider_gosdk_image_addr\":\"$(getFullImage sdkat_provider_gosdk)\",
    \"consumer_mesher_app_image_addr\":\"$(getFullImage sdkat_consumer_mesher)\",
    \"provider_mesher_app_image_addr\":\"$(getFullImage sdkat_provider_mesher)\",
    \"mesher_image_addr\":\"$MesherImage\"
  }
}"
fi

if [ "$1" = "gosdk" ];then
    curl -k -X PUT -H "$TokenHeaderKey:$tmpToken" "$Scheme$AosAddr$StackLifeApi" -d "{
  \"lifecycle\":\"upgrade\",
  \"inputs\":{
    \"consumer_gosdk_image_addr\":\"$(getFullImage sdkat_consumer_gosdk)\",
    \"provider_gosdk_image_addr\":\"$(getFullImage sdkat_provider_gosdk)\",
  }
}"
fi
if [ "$1" = "mesher" ];then
    curl -k -X PUT -H "$TokenHeaderKey:$tmpToken" "$Scheme$AosAddr$StackLifeApi" -d "{
  \"lifecycle\":\"upgrade\",
  \"inputs\":{
    \"consumer_mesher_app_image_addr\":\"$(getFullImage sdkat_consumer_mesher)\",
    \"provider_mesher_app_image_addr\":\"$(getFullImage sdkat_provider_mesher)\",
    \"mesher_image_addr\":\"$MesherImage\"
  }
}"
fi
