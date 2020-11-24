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

## Environment Configuration

The following is a list of environment variables to configure Krane
	
**KRANE_PRIVATE_KEY**

The private key used by Krane for signing authentication requests.

- required: `true`

**SECURED**

Enable HTTPS/TLS communication

- required: `false`

- default: `false`

**LISTEN_ADDRESS**

- default: `127.0.0.1:8500`

**KRANE_LOG_LEVEL**

- default: `info`

- values: `debug|info|warn|error`

**DB_PATH**

Krane uses [boltdb](https://github.com/etcd-io/bbolt) as its backing store for storing configuration details. Boltdb is represented as a single file on your disk, this is the path Krane will use to create/manage boltdb. Companies such as Shopify and Heroku use bolt within high-load production environments every day. 

- default: `/tmp/krane.db`   