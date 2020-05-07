#!/bin/bash

# Set local env
app=krane-server
owner=biensupernice
server_path=$PWD/cmd/server
bin_path=usr/local/bin # Need to check if this should be overridable..

echo "Installing $app"

# Verify docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "docker required to bootstrap $app"
    exit 0
fi

# Configure host 'on-start processes' file

# Start Server
execute(){
    "$bin_path/$app"
}

execute



