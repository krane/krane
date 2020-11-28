[![Krane](https://res.cloudinary.com/biensupernice/image/upload/v1602474802/Marketing_-_Krane_dj2y9e.png)](https://krane.sh)

[![CI](https://github.com/biensupernice/krane/workflows/CI/badge.svg?branch=master)](https://github.com/biensupernice/krane/actions)
[![Docker Pulls](https://img.shields.io/docker/pulls/biensupernice/krane?label=Docker%20Pulls)](https://store.docker.com/community/images/biensupernice/krane)
[![Go Report Card](https://goreportcard.com/badge/github.com/biensupernice/krane)](https://goreportcard.com/report/github.com/biensupernice/krane)

> ⚠️ Currently under construction

Krane is a self-hosted container management solution. It offers a simple and productive way to work with docker containers. Krane lets you deploy containers with a single configuration file and manage those containers using the Krane cli.

* **Documentation:** https://krane.sh  
* **Releases:** https://github.com/biensupernice/krane/releases
* **CLI:** https://github.com/krane/cli

## Features

* Single file deployments
* Provides HTTPS/TLS to your containers via [Let's Encrypt](https://letsencrypt.org/) 
* Deployment secrets
* Deployment scaling w/ container discovery 
* Round Robin load-balancing provided by [Traefik](https://doc.traefik.io/traefik/routing/services/#load-balancing)

## Getting Started

![Install Krane](./docs/assets/1-install-krane.png)

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY='changeme' \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -p 8500:8500 biensupernice/krane
```

Other [installation](https://www.krane.sh/#/installation) methods and configurations.


2. Download CLI

Download the Krane [CLI](https://www.krane.sh/#/cli) to execute commands on a Krane instance.

```
npm i -g @krane/cli
```

Full list of [commands](https://www.krane.sh/#/cli?id=commands).

3. Setup Authentication

Create public and private keys used for authentication.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'
```

The private key is kept on the user's machine, the public key is stored where Krane is running and appended to `~/.ssh/authorized_keys`

4. Authenticate

When logging in, you'll be prompted for the endpoint where Krane is running and the public key you created in step 3. Once authenticated successfully you'll be able to execute any command on that Krane instance. 

To switch between Krane instances you'll have to login again.

```
krane login
```

5. Deploy

Create a deployment configuration file `deployment.json` 

For example:

```
{
  "name": "hello-world-app",
  "image": "hello-world",
  "alias": ["hello.example.com"]
}
```

```
krane deploy -f /path/to/deployment.json
```

For more deployment configuration options, checkout the [documentation](https://www.krane.sh/#/deployment-config).

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
