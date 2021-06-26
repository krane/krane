#!/bin/bash

# -- check if the krane container is running --
is_krane_running() {
  local container_name="krane"
  if [ "$(docker container inspect -f '{{.State.Status}}' "$container_name")" == "running" ];
    then
      return 0
    else
      return 1
  fi
}

# -- fatal if the system does not meet minimum requirements --
verify_system() {
  # -- check docker is installed --
  if ! [ -x "$(command -v docker)" ]; then
    echo "Docker required to properly install Krane"
    exit 0
  fi

  # -- check docker is running --
  if ! docker info >/dev/null 2>&1; then
    echo "Docker must be in a RUNNING state to install Krane"
    echo "try running [$ docker info] to troubleshoot your Docker daemon"
    exit 0
  fi

  if is_krane_running ; then
    prompt_yes_no EXIT_EARLY "A Krane instance is already running. Would you like to continue?"
  fi

  if [ "$EXIT_EARLY" == true ]; then
    exit 0
  fi
}

# -- setup the system environment --
setup_env() {
  echo "Installing Krane "
  echo -e "\nThis interactive script will configure:"
  echo "• A running Krane container instance"
  echo "• A running proxy container to route DNS traffic"
  echo -e "\nPlease see the official docs site for complete documentation:\nhttps://krane.sh/#/docs/installation \n"

  prompt_yes_no IS_LOCAL "Are you running Krane on localhost?"

  if [ "$IS_LOCAL" == true ];
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
    ensure_env LETSENCRYPT_EMAIL "What email should we use to generate deployment TLS certificates?"
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
}

# --helper function to ensure environment variable is set --
ensure_env() {
  local env=$1
  local prompt_msg=$2

  if [ ! -z "$env" ]; then
    read -p "$prompt_msg: `echo $'\n$ '`" value
    export "$env=$value"
  fi
}

# -- helper function to ensure sensitive environment variable is set --
ensure_secure_env() {
  local env=$1
  local prompt_msg=$2

  if [ ! -z "$env" ]; then
    read -s -p "$prompt_msg: `echo $'\n$ 🔒'`" value
    export "$env=$value"
  fi
  echo ""
}

# -- prompt the console for yes or no response --
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

# -- download Krane and verify in it's in a running state --
download_and_verify() {
  echo -e "✅ Stopping Krane (if running)"
  docker stop krane > /dev/null 2>&1

  echo -e "✅ Removing Krane instance (if exists)"
  docker rm krane > /dev/null 2>&1

  echo -e "✅ Removing existing image"
  docker image rm biensupernice/krane > /dev/null 2>&1

  echo -e "✅ Pulling latest image"
  docker pull biensupernice/krane:latest -q > /dev/null 2>&1

  echo -e "✅ Preparing Krane network"
  docker network create --driver bridge krane > /dev/null 2>&1

  echo -e "✅ Starting new Krane instance \n"
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

  echo -e "\n⏳ Waiting 10 seconds for Krane to be ready!"
  sleep 10s

  # -- check that the Krane container is in a running state --
  local container_name="krane"
  if ! [ "$(docker container inspect -f '{{.State.Status}}' $container_name)" == "running" ];  then
    echo -e "\n❌ Encountered an error starting Krane:"
    docker logs krane 2>&1
    exit 0
  fi

  echo -e "\n✅ Cleaning up older images"
  docker image prune -a -f

  echo -e "\nInstallation complete..."
  echo -e "For documentation on accessing this Krane instance visit:\n https://www.krane.sh/#/docs/cli"
  
  echo -e "\nTake note of the following details used to create your Krane instance:"
  echo -e "🔐 Krane private key: $KRANE_PRIVATE_KEY"
  echo -e "📝 SSH keys directory: $SSH_KEYS_DIR"
  echo -e "🗂 Database path: $DB_PATH"
  echo -e "🌎 Proxy alias: $PROXY_DASHBOARD_ALIAS"

  echo -e "\nSome helpful commands to help you start using this Krane instance:"
  echo -e "$ krane login http://$ROOT_DOMAIN:8500"
  echo -e "$ krane ls"
  echo -e "$ krane deploy -f ./deployment.json"
}

# --- run the install process --
{
  verify_system
  setup_env
  download_and_verify
}