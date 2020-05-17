#!/bin/bash

set -e

APP=krane-server
SERVER_ENTRY_PATH=$PWD/cmd/krane-server
BIN_PATH=/usr/local/bin

if [ "$1" = 'build' ]; then
    # Verify go is installed
    if [ ! -x "$(command -v go)" ]; then
        echo "go required to package $app"
        exit 0
    fi

    # Run the build
    go build -o $BIN_PATH/$APP $SERVER_ENTRY_PATH
    echo "> $APP installed @ $BIN_PATH/$APP"
    echo "> $APP # run the service"
fi

