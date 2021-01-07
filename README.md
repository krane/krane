[![Krane](https://res.cloudinary.com/biensupernice/image/upload/v1602474802/Marketing_-_Krane_dj2y9e.png)](https://krane.sh)

[![CI](https://github.com/krane/krane/workflows/CI/badge.svg?branch=master)](https://github.com/krane/krane/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/krane/krane)](https://goreportcard.com/report/github.com/krane/krane)
[![Coverage](https://img.shields.io/codecov/c/github/krane/krane?color=blue)](https://codecov.io/gh/krane/krane)

> ⚠️ Currently under construction

Krane makes it easy to deploy containers for development workloads on remote or local servers. Krane interfaces with Docker exposing a productive toolset for managing containerized services known as deployments. The Krane [CLI](https://www.krane.sh/#/cli) allows you to interact with Krane to create, manage and automate deployments.

- **Documentation:** https://krane.sh
- **Releases:** https://github.com/krane/krane/releases
- **CLI:** https://github.com/krane/cli
- **GitHub Action:** https://github.com/krane/action

## Features

- Single file deployments
- Compatible with _localhost_ with features like aliases(`my-api.localhost`)
- HTTPS/TLS support via [Let's Encrypt](https://letsencrypt.org/)
- Deployment [aliases](https://www.krane.sh/#/deployment-configuration?id=alias) provided by [Traefik](https://traefik.io/traefik/)
- Deployment [secrets](https://www.krane.sh/#/deployment-configuration?id=secrets) for hiding sensitive environment variables
- Deployment [scaling](https://www.krane.sh/#/deployment-configuration?id=scale) to distribute the workload between containers
- Tooling - [CLI](https://github.com/krane/cli), [GitHub Action](https://github.com/krane/action)
- [Self-hosted](#motivation) - Bring your own server (could be a cheap $5 server) and scale if you need

## Getting Started

[![Install Krane](website/assets/1-install-krane.png)](https://www.krane.sh/#/installation)

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY=changeme \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh:/root/.ssh \
    -p 8500:8500 biensupernice/krane
```

Other [installation](https://www.krane.sh/#/installation) methods and configurations.

[![Download CLI](website/assets/2-download-cli.png)](https://www.krane.sh/#/cli)

Download the Krane [CLI](https://www.krane.sh/#/cli) to execute commands on a Krane instance.

```
npm i -g @krane/cli
```

Full list of [commands](https://www.krane.sh/#/cli?id=commands).

![Setup Authentication](website/assets/3-setup-authentication.png)

Create public and private keys used for authentication.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'
```

The private key stays on the user's machine, the public key is appended to `~/.ssh/authorized_keys` where Krane is running.

[![Authenticate](website/assets/4-authentication.png)](https://www.krane.sh/#/cli?id=authenticating)

When logging in, you'll be prompted for the endpoint where Krane is running and the public key you created in step 3. Once authenticated you'll be able to execute commands on that Krane instance.

To switch between Krane instances you'll have to login again.

```
krane login
```

[![Deploy](website/assets/5-deploy.png)](https://www.krane.sh/#/cli?id=deploy)

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

Krane is a self-hosted SaaS container tool. You bring your own server and install Krane on it to manage your containers in the form of deployments - The benefit is _cost per deployment_. Pricing of other platforms such as Digital Ocean's [app-platform](https://www.digitalocean.com/docs/app-platform/) start at $5 per deployment. A self-hosted solution allows you to own your server (cheap), and the benefit of multiple deployments for no extra cost. Maintaining and managing your own solution may sound complex, Krane tries to make the process straight-forward and cost-effective.

Krane isn't a replacement for [Kubernetes](https://kubernetes.io/), [ECS](https://aws.amazon.com/ecs/), or any other container orchestration solution you might see running production applications, instead it's a tool you can leverage to make development of side-projects or small workloads cheap and straight forward. In the end, that was the main objective, a productive deployment tool for managing non-critical container workloads on remote servers.

## Building from source

```
$ git clone https://github.com/krane/krane
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
