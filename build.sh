#!/bin/bash

app=krane-server
owner=biensupernice
server_path="$PWD"/cmd/krane-server
bin_path=/usr/local/bin # Need to check if this should be overridable..

echo "> Building $app"

# Verify go is installed
if [ ! -x "$(command -v go)" ]; then
    echo "go required to package $app"
    exit 0
fi

# Run the build
cd "$server_path"

go build -o "$bin_path"/$app -tags $owner
echo "> $app installed succesfully"
echo $'\nRun:'
echo "> $app"