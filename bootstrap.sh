#!/bin/bash

# Verify docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker required to bootstrap krane-server"
    exit 0
fi


# Set default env
KRANE_SERVER_PATH=""

# Configure host 'on-start processes' file

# Start Server
$KRANE_SERVER_PATH/server

