# Krane Config

Deploying containers using **Krane** starts with a configuration file that describes how Krane should run your containers.

> A recommended pattern is to have a `krane.json` file at the root of you project.

An example of a `krane.json` for deploying the Krane UI.

```json
{
  "name": "krane-ui",
  "image": "biensupernice/krane-ui",
  "env": {
    "NODE_ENV": "dev"
  },
  "secrets": {
    "KRANE_HOST": "@krane-host",
    "KRANE_TOKEN": "@krane-token"
  }
}
```

## name

The name of your deployment.

- required: `yes`

## registry

The container registry to use when pulling images.

- required: `false`
- default: `docker.io`

## image

The image used when pulling, creating and running your deployments containers.

- required: `true`

## env

The enviornment variables passed to the containers part of a deployment.

> ⚠️ Enviornment variables should not contain any sensitive data, use `secrets` instead.

- required: `false`

```json
{
  "env": {
    "NODE_ENV": "dev",
    "PORT": "8080"
  }
}
```

## secrets

Secrets are used when you want to pass sensitive information to your deployment. Secrets are **not shared** across deployments, they are bounded and only provided to the containers part of the deployment.

Secrets are created using the krane `cli` and referenced in your Krane configuration using the `@` symbol.

- required: `false`

```json
{
  "secrets": {
    "SECRET_TOKEN": "@my-secret-token",
    "PROXY_API_URL": "@proxy-api-url"
  }
}
```

## tag

The tag used when pulling the image.

- required: `false`
- default: `latest`

## volumes

The volumes to mount from the container to the host.

- required: `false`

```json
{
  "volumes": {
    "/home/user/data": "/data/db"
  }
}
```
