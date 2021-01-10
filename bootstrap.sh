#!/bin/bash

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

update_or_create_krane(){
  echo "(1/6) Stopping Krane (if exists)"
  docker stop krane

  echo -e "\n(2/6) Removing Krane instance (if exists)"
  docker rm krane

  echo -e "\n(3/6) Removing existing image"
  docker image rm biensupernice/krane

  echo -e "\n(4/6) Pulling latest image"
  docker pull biensupernice/krane:latest -q

  echo -e "\n(5/6) Starting new Krane instance"
  docker run -d --name=krane \
    -e LOG_LEVEL=info \
    -e KRANE_PRIVATE_KEY="${KRANE_PRIVATE_KEY:-uuidgen}" \
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

  echo -e "\n(6/6) Cleaning up older images"
  docker image prune -a -f

  echo -e "\nBootstrap complete."
  echo -e "\nNote: Make sure no errors where found above before attempting to access your Krane instance."
  echo "For documentation on accessing this Krane instance visit https://www.krane.sh/#/docs/cli."
}

ensure_env KRANE_PRIVATE_KEY "Krane private key (optional, used for signing client requests. default uuid)"
ensure_env DOCKER_BASIC_AUTH_USERNAME "Container registry username (optional, will operate as an anonymous user)"
ensure_secure_env DOCKER_BASIC_AUTH_PASSWORD "Container registry password (optional, will operate as an anonymous user)"
ensure_env SSH_KEYS_DIR "SSH keys directory (optional, default /root/.ssh)"
ensure_env DB_DIR "Krane database directory (optional, default /tmp)"
ensure_env PROXY_ENABLED "Network proxy enabled? (default true, Note: if aliases are needed set as true)"
ensure_env PROXY_DASHBOARD_ALIAS "Network proxy dashboard alias (domain alias for the proxy dashboard, for example monitor.example.com)"
ensure_env PROXY_DASHBOARD_SECURE "Network proxy secure? (default true, enables https on the network proxy. Note: should be set to false for 'localhost')"
ensure_env LETSENCRYPT_EMAIL "Certificate email: (optional, email used to generate https/tls certificates. Note: Must be a valid email for Let's Encrypt to properly generate certs. Ignore for 'localhost')"

update_or_create_krane