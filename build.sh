#!/bin/bash

APP=krane-server
OWNER=biensupernice
SERVER_ENTRY_PATH=$PWD/cmd/krane-server
BIN_PATH=/usr/local/bin

echo "> Building $APP in $SERVER_ENTRY_PATH"

if [ "$1" = 'start' ]; then
    # Verify go is installed
    if [ ! -x "$(command -v go)" ]; then
        echo "go required to package $app"
        exit 0
    fi

    # Run the build
    cd "$SERVER_ENTRY_PATH"
    go build -o "$BIN_PATH"/$APP -tags $OWNER

    echo "🏗 Starting $APP"

    export KRANE_PRIVATE_KEY=${PRIVATE_KEY:-"2733dd1ccfe8d36b6ec8818c78a8940ee714237f"}
    export KRANE_PORT=${PORT:-8080}
    export KRANE_LOG_LEVEL=${LOG_LEVEL:-"debug"}
    export KRANE_PATH=${KRANE_PATH:-"/.krane"}

    echo "> $APP port: $KRANE_PORT"
    echo "> $APP log level: $KRANE_LOG_LEVEL"
    echo "> $APP path: $KRANE_PATH"

    mkdir $KRANE_PATH

    echo "\n> $APP installed succesfully"

    sh -c krane-server
fi

