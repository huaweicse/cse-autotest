#!/bin/bash
set -x

docker ps -a|grep sdkat|awk '{print $1}'|xargs docker rm -f
