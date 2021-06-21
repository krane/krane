<!-- [![Logo](docs/assets/logo.png)](https://krane.sh) -->

# Krane

> Open-source, self-hosted, container management solution

[![CI](https://github.com/krane/krane/workflows/CI/badge.svg?branch=master)](https://github.com/krane/krane/actions)
[![Release](https://img.shields.io/github/v/release/krane/krane)](https://github.com/krane/krane/releases)

Krane is a container management solution that abstracts away the hard parts from deploying infrastructrue at the lowest cost possible.

- **Documentation:** https://krane.sh
- **Releases:** https://github.com/krane/krane/releases
- **Bugs:** https://github.com/krane/krane/issues

## Krane Tooling

- **Deployment Dashboard:** https://github.com/krane/ui
- **CLI:** https://github.com/krane/cli
- **GitHub Action:** https://github.com/krane/action

## Features

- One command deployments
- Single file deployments
- Deployment [aliases](https://www.krane.sh/#/docs/deployment?id=alias) (`my-api.localhost`)
- Deployment [secrets](https://www.krane.sh/#/docs/deployment?id=secrets) for hiding sensitive environment variables
- Deployment [scaling](https://www.krane.sh/#/docs/deployment?id=scale) to distribute the workload between containers
- Deployment [rate limit](https://www.krane.sh/#/docs/deployment?id=rate_limit) to limit incoming requests
- HTTPS/TLS support with auto generated [Let's Encrypt](https://letsencrypt.org/) certificates
- [Self-hosted](#motivation) - Cost-effective, bring your own server and scale if you need

## Getting Started

1. Install Krane

```
bash <(wget -qO- get.krane.sh)
```

2. Create a deployment configuration file

`deployment.json`

```json
{
  "name": "krane-getting-started",
  "image": "docker/getting-started",
  "alias": ["getting-started.example.com"]
}
```

3. Deploy your infrastructure

```
krane deploy -f ./deployment.json
```

For more deployment configuration options, checkout the [documentation](https://www.krane.sh/#/docs/deployment)

<a name="motivation"></a>

## Motivation

Krane is a self-hosted PaaS. You bring your own server and install Krane on it to manage your containers in the form of deployments - The benefit, <i>cost per deployment</i>. A self-hosted solution allows you to own your server (cost-effective), and the benefit of any number of deployments at no extra cost. Maintaining and managing your own solution may sound complex, Krane tries to make the process <i>straight-forward</i> and <i>cost-effective</i> .

Krane isn't a replacement for [Kubernetes](https://kubernetes.io), [ECS](https://aws.amazon.com/ecs/), or any other container orchestration solution you might see running production applications, instead it's a tool you can leverage to make development of side-projects or small workloads cheap and straight forward. That was the main objective, a productive deployment tool for managing non-critical container workloads on remote servers.

## Building from source

```
$ git clone https://github.com/krane/krane
$ cd krane
$ go build ./cmd/krane
$ export KRANE_PRIVATE_KEY=changeme
$ ./krane
```

## Running tests

[![Go Report Card](https://goreportcard.com/badge/github.com/krane/krane)](https://goreportcard.com/report/github.com/krane/krane)
[![Coverage](https://img.shields.io/codecov/c/github/krane/krane?color=blue)](https://codecov.io/gh/krane/krane)

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

## Minimal Docker example

This is the most minimal Docker example to get _up-and-running_ with Krane

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY=changeme \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh:/root/.ssh  \
    -p 8500:8500 biensupernice/krane
```

## Complete Docker example

This is a complete Docker example to get Krane running with:

- Automatic HTTPS/SSL w/ Lets Encrypt certificates
- Container registry authentication for pulling images
- Volumed Krane DB (for storing session & deployment details)
- Log level set to debug (for debugging)

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
