#!/bin/sh

CURRENT_DIR=$(cd $(dirname $0);pwd)
cd $CURRENT_DIR

LOG() {
    if [ "$1" != "" ]; then
        echo $(date "+%Y-%m-%d %H:%M:%S") $1
    fi
}

LOG "Run test case after 30 seconds"
sleep 30
go test ./... -v

echo "
                            ######## ##    ## ########
                            ##       ###   ## ##     ##
                            ##       ####  ## ##     ##
                            ######   ## ## ## ##     ##
                            ##       ##  #### ##     ##
                            ##       ##   ### ##     ##
                            ######## ##    ## ########
                            "
sleep 30000000
