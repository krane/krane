#!/bin/bash

# Check if docker is installed
if ! [ -x "$(command -v docker)" ]; then
  echo "Docker required to properly bootstrap Krane"
  exit 0
fi

# Check if docker is running
if ! docker info >/dev/null 2>&1; then 
  echo "Docker must be in a RUNNING state to bootstrap Krane"
  echo "try running [$ docker info] to troubleshoot your Docker daemon"
  exit 0
fi

ensure_env() {
  local env=$1
  local prompt_msg=$2

  if [ ! -z "$env" ]; then
    read -p "$prompt_msg: `echo $'\n$ '`" value
    export "$env=$value"
  fi
}

ensure_secure_env() {
  local env=$1
  local prompt_msg=$2

  if [ ! -z "$env" ]; then
    read -s -p "$prompt_msg: `echo $'\n$ 🔒'`" value
    export "$env=$value"
  fi
  echo ""
}

prompt_yes_no() {
  local env=$1
  local prompt_msg=$2

  while true; do
    read -p "$prompt_msg yes [Yy] / no [Nn]`echo $'\n$ '`" yn
    case $yn in
        [Yy]* ) export "$env=true"; break;;
        [Nn]* ) export "$env=false"; break;;
        * ) echo "Please answer yes[Yy] or no[Nn].";;
    esac
done
}

bootstrap_krane() {
  echo -e "(1/7) Stopping Krane (if exists)"
  docker stop krane > /dev/null 2>&1

  echo -e "(2/7) Removing Krane instance (if exists)"
  docker rm krane > /dev/null 2>&1

  echo -e "(3/7) Removing existing image"
  docker image rm biensupernice/krane > /dev/null 2>&1

  echo -e "(4/7) Pulling latest image"
  docker pull biensupernice/krane:latest -q > /dev/null 2>&1

  echo -e "(5/7) Ensuring Krane network"
  docker network create --driver bridge krane > /dev/null 2>&1

  echo -e "(6/7) Starting new Krane instance \n"
  docker run -d --name=krane --network=krane \
    -e LOG_LEVEL=info \
    -e KRANE_PRIVATE_KEY="${KRANE_PRIVATE_KEY:-$(uuidgen)}" \
    -e DB_PATH=${DB_PATH} \
    -e DOCKER_BASIC_AUTH_USERNAME="$DOCKER_BASIC_AUTH_USERNAME" \
    -e DOCKER_BASIC_AUTH_PASSWORD="$DOCKER_BASIC_AUTH_PASSWORD" \
    -e PROXY_ENABLED="${PROXY_ENABLED:-true}" \
    -e PROXY_DASHBOARD_ALIAS="$PROXY_DASHBOARD_ALIAS" \
    -e PROXY_DASHBOARD_SECURE="${PROXY_DASHBOARD_SECURE:-true}" \
    -e LETSENCRYPT_EMAIL="$LETSENCRYPT_EMAIL" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v "${SSH_KEYS_DIR:-/root/.ssh}":/root/.ssh  \
    -v "${DB_DIR:-/tmp}":/tmp \
    -p 8500:8500 biensupernice/krane

  sleep 5s

  # Check that the Krane containers is in a running state 
  local container_name="krane"
  if ! [ "$(docker container inspect -f '{{.State.Status}}' $container_name)" == "running" ];  then
    echo -e "\nEncountered and error while bootstrapping Krane."
    exit 0
  fi

  echo -e "\n(7/7) Cleaning up older images"
  docker image prune -a -f

  echo -e "\nBootstrap complete..."
  echo -e "For documentation on accessing this Krane instance visit https://www.krane.sh/#/docs/cli"
  
  echo -e "\nTake note of the following details used to create your Krane instance:"
  echo -e "* Krane private key: $KRANE_PRIVATE_KEY"
  echo -e "* SSH keys directory: $SSH_KEYS_DIR"
  echo -e "* Database directory: $DB_DIR"
  echo -e "* Proxy alias: $PROXY_DASHBOARD_ALIAS"

  echo -e "\nHere are some helpful commands to help you start using this Krane instance:"
  echo -e "$ krane login http://$ROOT_DOMAIN:8500"
  echo -e "$ krane ls"
  echo -e "$ krane deploy -f ./deployment.json"
}

echo "Bootstrapping Krane..."
echo -e "\nThis interactive script will help you setup:"
echo "• Krane instance to manage containers"
echo "• Krane proxy to route traffic"
echo -e "\nFor complete documentation visit https://krane.sh/#/docs/installation \n"

prompt_yes_no IS_LOCAL "Are you running Krane on localhost?"

if [ $IS_LOCAL ];
then
  export ROOT_DOMAIN="localhost"
  export KRANE_PRIVATE_KEY="krane"  
  export SSH_KEYS_DIR="$HOME/.ssh"
  export DB_PATH="/tmp/krane.db"
  export PROXY_ENABLED=true
  export PROXY_DASHBOARD_ALIAS="proxy.$ROOT_DOMAIN"
  export PROXY_DASHBOARD_SECURE=false
else
  ensure_env ROOT_DOMAIN "What domain do you want to use for this Krane instance? (ie. example.com | krane.example.com)"
  ensure_env LETSENCRYPT_EMAIL "What email should we use to generate your deployments TLS certificates?"
  ensure_env KRANE_PRIVATE_KEY "What should we use as the Krane private key (used for signing client requests, defaults to a uuid)"
  ensure_env SSH_KEYS_DIR "What directory are your SSH keys located? (optional, default directory /root/.ssh)"
  ensure_env DB_DIR "What directory should we use for the Krane database? (optional, default directory /tmp)"
  ensure_env DOCKER_BASIC_AUTH_USERNAME "What is the container registry username you want to use? (optional, will operate as an anonymous user)"
  ensure_secure_env DOCKER_BASIC_AUTH_PASSWORD "What is the container registry password? (optional, will operate as an anonymous user)"
  export DB_PATH="${DB_DIR:-/tmp}/krane.db"
  export PROXY_ENABLED=true
  export PROXY_DASHBOARD_ALIAS="proxy.$ROOT_DOMAIN"
  export PROXY_DASHBOARD_SECURE=true
fi

bootstrap_krane