#!/bin/bash

DEFAULT_KRANE_PRIVATE_KEY="1ab43c39-2cec-4a0a-b088-67a4361c5714"

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
    prompt_yes_no CONTINUE_INSTALLATION "A Krane instance is already running. Would you like to continue?"
  else
    export CONTINUE_INSTALLATION="true"
  fi

  if ! [ "$CONTINUE_INSTALLATION" == true ]; then
    exit 0
  fi
}

# -- setup the system environment --
setup_env() {
  echo -e "\nInstalling Krane..."
  echo -e "\nThis interactive script will configure:"
  echo "• A running Krane container instance"
  echo "• A running Krane proxy container to route DNS traffic"
  echo -e "\nCheckout the official docs site for complete documentation:\nhttps://krane.sh/#/docs/installation \n"

  prompt_yes_no IS_LOCAL "Are you running Krane on localhost?"

  if [ "$IS_LOCAL" == true ];
  then
    export ROOT_DOMAIN="krane.localhost"
    export SSH_KEYS_DIR="$HOME/.ssh"
    export DB_PATH="/tmp/krane.db"
    export PROXY_ENABLED=true
    export PROXY_DASHBOARD_ALIAS="proxy.$ROOT_DOMAIN"
    export PROXY_DASHBOARD_SECURE=false
  else
    ensure_env ROOT_DOMAIN "What domain do you want to use for this Krane instance? (ie. example.com | krane.example.com)"
    ensure_env LETSENCRYPT_EMAIL "What email should we use to generate deployment TLS certificates?"
    ensure_env SSH_KEYS_DIR "What directory are your SSH keys located? (optional, default directory /root/.ssh)"
    ensure_env DB_DIR "What directory should we use for the Krane database? (optional, default directory /tmp)"
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
  echo -e "\n✓ Stopping Krane (if running)"
  docker stop krane > /dev/null 2>&1

  echo -e "✓ Removing Krane instance (if exists)"
  docker rm krane > /dev/null 2>&1

  echo -e "✓ Pulling latest image"
  docker image rm biensupernice/krane > /dev/null 2>&1
  docker pull biensupernice/krane:latest -q > /dev/null 2>&1

  echo -e "✓ Preparing Krane network"
  docker network create --driver bridge krane > /dev/null 2>&1

  echo -e "✓ Starting new Krane instance"
  docker run -d --name=krane --network=krane \
    -e LOG_LEVEL=info \
    -e KRANE_PRIVATE_KEY="${KRANE_PRIVATE_KEY:-${DEFAULT_KRANE_PRIVATE_KEY}}" \
    -e DB_PATH=${DB_PATH} \
    -e PROXY_ENABLED="${PROXY_ENABLED:-true}" \
    -e PROXY_DASHBOARD_ALIAS="$PROXY_DASHBOARD_ALIAS" \
    -e PROXY_DASHBOARD_SECURE="${PROXY_DASHBOARD_SECURE:-true}" \
    -e LETSENCRYPT_EMAIL="$LETSENCRYPT_EMAIL" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v "${SSH_KEYS_DIR:-/root/.ssh}":/root/.ssh  \
    -v "${DB_DIR:-/tmp}":"${DB_DIR:-/tmp}" \
    -l "traefik.enable=true" \
    -l "traefik.http.middlewares.redirect-to-https.redirectscheme.permanent=true" \
    -l "traefik.http.middlewares.redirect-to-https.redirectscheme.port=443" \
    -l "traefik.http.middlewares.redirect-to-https.redirectscheme.scheme=https" \
    -l "traefik.http.routers.krane-insecure.entrypoints=web" \
    -l "traefik.http.routers.krane-insecure.middlewares=redirect-to-https" \
    -l "traefik.http.routers.krane-insecure.rule=Host(\`${ROOT_DOMAIN}\`)" \
    -l "traefik.http.routers.krane-secure.entrypoints=web-secure" \
    -l "traefik.http.routers.krane-secure.rule=Host(\`${ROOT_DOMAIN}\`)" \
    -l "traefik.http.routers.krane-secure.tls=true" \
    -l "traefik.http.routers.krane-secure.tls.certresolver=lets-encrypt" \
    -p 8500:8500 biensupernice/krane > /dev/null 2>&1

  echo -e "\n⏳ Waiting for Krane to be ready"
  sleep 10s

  # -- check that the Krane container is in a running state --
  local container_name="krane"
  if ! [ "$(docker container inspect -f '{{.State.Status}}' $container_name)" == "running" ];  then
    echo -e "\n✕ Encountered an error starting Krane:"
    docker logs krane 2>&1
    exit 0
  fi

  echo -e "\n✓ Cleaning up older images"
  docker image prune -a -f

  echo -e "\n✓ Installation complete 💛"

  echo -e "\nYou can now use your Krane instance with:"
  echo -e "🔐 Krane private key: ${KRANE_PRIVATE_KEY:-${DEFAULT_KRANE_PRIVATE_KEY}}"
  echo -e "📝 SSH keys directory: ${SSH_KEYS_DIR:-/root/.ssh}"
  echo -e "🗂  Database path: $DB_PATH"
  echo -e "🌎 Proxy alias: $PROXY_DASHBOARD_ALIAS"
  echo -e "🕹  Let's Encrypt email: ${LETSENCRYPT_EMAIL:-not set}"

  echo -e "\nSome helpful CLI commands to get you started:"
  echo -e "$ krane login http://$ROOT_DOMAIN"
  echo -e "$ krane ls"
  echo -e "$ krane deploy -f ./deployment.json"

  echo -e "\nThanks for using Krane! ☺️"
}

# --- run the install process --
{
  verify_system
  setup_env
  download_and_verify
}
