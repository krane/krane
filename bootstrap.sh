#!/bin/bash

ensure_env() {
  local env=$1
  local prompt_msg=$2

  if [ ! -z "$env" ]; then
    read -p "$prompt_msg `echo $'\n> '`" value
    export "$env=$value"
  fi
  printf "\n"
}

ensure_secure_env() {
  local env=$1
  local prompt_msg=$2

  if [ ! -z "$env" ]; then
    read -s -p "$prompt_msg `echo $'\n> '`" value
    export "$env=$value"
    printf "\n"
  fi
  printf "\n"
}

# steps to create a Krane instance using docker
create_krane_docker(){
  echo "Stopping Krane"
  docker stop krane

  echo -e "\nRemoving Krane instance"
  docker rm krane

  echo -e "\nRemoving existing image"
  docker image rm biensupernice/krane

  echo -e "\nPulling latest image"
  docker pull biensupernice/krane:latest -q

  echo -e "\nStarting new instance"
  docker run -d --name=krane \
    -e LOG_LEVEL=info \
    -e KRANE_PRIVATE_KEY="$KRANE_PRIVATE_KEY" \
    -e DB_PATH="${DB_DIR/krane.db:-/tmp/krane.db}" \
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

  echo -e "\nCleaning up older images"
  docker image prune -a -f
}

ensure_env DOCKER_BASIC_AUTH_USERNAME "Container registry username: (optional, will operate as an anonymous user. Note: May affect registry rate-limits)"
ensure_secure_env DOCKER_BASIC_AUTH_PASSWORD "Container registry password: (optional, only if container registry username is being used)"
ensure_env SSH_KEYS_DIR "Ssh keys directory: (optional, default /root/.ssh)"
ensure_env DB_DIR "Krane database directory: (optional, default /tmp)"
ensure_secure_env KRANE_PRIVATE_KEY "Krane private key: (required, for signing client requests)"
ensure_env PROXY_ENABLED "Network proxy enabled? (default true, Note: if aliases are needed set as true)"
ensure_env PROXY_DASHBOARD_ALIAS "Network proxy dashboard alias: (alias for the proxy dashboard, ie. monitor.example.com)"
ensure_env PROXY_DASHBOARD_SECURE "Network proxy secure? (default true, enable https on the network proxy)"
ensure_env LETSENCRYPT_EMAIL "Certificate email: (optional, email used to generate https/tls certificates. Note: Must be valid for let's encrypt to generate certs)"

create_krane_docker