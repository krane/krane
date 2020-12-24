[![Krane](https://res.cloudinary.com/biensupernice/image/upload/v1602474802/Marketing_-_Krane_dj2y9e.png)](https://krane.sh)

[![CI](https://github.com/biensupernice/krane/workflows/CI/badge.svg?branch=master)](https://github.com/biensupernice/krane/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/biensupernice/krane)](https://goreportcard.com/report/github.com/biensupernice/krane)
[![Coverage](https://img.shields.io/codecov/c/github/biensupernice/krane?color=blue)](https://codecov.io/gh/biensupernice/krane)

> ⚠️ Currently under construction

Krane is a container deployment service that makes it easy to create and manage small development workloads. Krane sits on any server and interfaces with Docker exposing a productive toolset for managing containers. The Krane [CLI](https://www.krane.sh/#/cli) allows you to automate or deploy application resources from anywhere.

- **Documentation:** https://krane.sh
- **Releases:** https://github.com/biensupernice/krane/releases
- **CLI:** https://github.com/krane/cli

## Features

- Single file deployments
- Provides HTTPS/TLS to your containers via [Let's Encrypt](https://letsencrypt.org/)
- Deployment [secrets](https://www.krane.sh/#/cli?id=secrets)
- Deployment [scaling](https://www.krane.sh/#/deployment-configuration?id=scale)
- Round Robin load-balancing provided by [Traefik](https://doc.traefik.io/traefik/routing/services/#load-balancing)
- [Self-hosted](#motivation) - Bring your own hardware (could be a cheap $5 machine) and scale if you need

## Getting Started

[![Install Krane](./docs/assets/1-install-krane.png)](https://www.krane.sh/#/installation)

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY='changeme' \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -p 8500:8500 biensupernice/krane
```

Other [installation](https://www.krane.sh/#/installation) methods and configurations.

[![Download CLI](./docs/assets/2-download-cli.png)](https://www.krane.sh/#/cli)

Download the Krane [CLI](https://www.krane.sh/#/cli) to execute commands on a Krane instance.

```
npm i -g @krane/cli
```

Full list of [commands](https://www.krane.sh/#/cli?id=commands).

![Setup Authentication](./docs/assets/3-setup-authentication.png)

Create public and private keys used for authentication.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'
```

The private key is kept on the user's machine, the public key is stored where Krane is running and appended to `~/.ssh/authorized_keys`

[![Authenticate](./docs/assets/4-authentication.png)](https://www.krane.sh/#/cli?id=authenticating)

When logging in, you'll be prompted for the endpoint where Krane is running and the public key you created in step 3. Once authenticated successfully you'll be able to execute any command on that Krane instance.

To switch between Krane instances you'll have to login again.

```
krane login
```

[![Deploy](./docs/assets/5-deploy.png)](https://www.krane.sh/#/cli?id=deploy)

Create a deployment configuration file `deployment.json`

For example:

```json
{
  "name": "hello-world-app",
  "image": "hello-world",
  "alias": ["hello.example.com"]
}
```

```
krane deploy -f /path/to/deployment.json
```

For more deployment configuration options, checkout the [documentation](https://www.krane.sh/#/deployment-configuration).

---

<a name="motivation"></a>

## Motivation

Krane is self-hosted, meaning you bring your own device and install Krane on it. The key benefit is **cost per deployment**. When comparing prices against other container deployment platforms, prices goes up pretty quickly when running multiple instances of your apps. Take Digital Ocean, deploying a container on [app-platform](https://www.digitalocean.com/docs/app-platform/) starts at $5 per instance, this could be the cost of a single Digital Ocean droplet running Krane but with many more deployments and functionality provided and built for you.

Krane isn't a replacement for Kubernetes or ECS or any other container management solution you might see running production applications, instead it's a tool you can leverage to make development of side-projects or small workloads cheap and straight forward. In the end, this was the main objective, a productive deployment tool for managing non-critical container workloads on remote servers.

## Building from source

```
$ git clone https://github.com/biensupernice/krane
$ cd krane
$ go build ./cmd/krane
$ export KRANE_PRIVATE_KEY=changeme
$ ./krane
```

## Running tests

In the root of the project

```
# run tests
$ go test ./...

or

# run tests with coverage
$ go test -coverprofile coverage.out ./...

# view coverage
$ go tool cover -html=coverage.out
```

## Viewing the database

Krane uses [boltdb](https://github.com/etcd-io/bbolt) as its backing store. To view the contents in bolt, you can use [boltdbweb](https://github.com/evnix/boltdbweb).

```
$ boltdbweb --db-name=/path/to/krane.db --port=9000
```

## Complete Docker example

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY=changeme \
    -e LOG_LEVEL=debug \
    -e DOCKER_BASIC_AUTH_USERNAME=changeme \
    -e DOCKER_BASIC_AUTH_PASSWORD=changeme \
    -e PROXY_ENABLED=true \
    -e PROXY_DASHBOARD_SECURE=true \
    -e PROXY_DASHBOARD_ALIAS=monitor.example.com \
    -e LETSENCRYPT_EMAIL=email@example.com \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh:/root/.ssh  \
    -v /tmp/krane.db:/tmp/krane.db \
    -p 8500:8500 biensupernice/krane
```
