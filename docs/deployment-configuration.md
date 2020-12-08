# Deployment Config

Creating a deployment using Krane starts with a **single file**, this file contains the deployment configuration used to create container resources. The deployment configuration can be stored anywhere, when using the [CLI](cli) you'll be referencing the location of this deployment configuration.

`deployment.json`

```json
{
  "name": "hello-world-app",
  "image": "hello-world",
  "alias": ["hello.example.com"]
}
```

The above **deployment configuration** sets up:

1. A deployment called **hello-world-app**
2. A container running the image **hello-world**
3. An alias **hello.example.com**

See all deployment configuration [options](deployment-configuration?id=options)

## Options

The different properties you can specificy in a deployment configuration file.

> A common pattern is to have a `deployment.json` at the root of your project

### name

The name of your deployment.

- required: `yes`

### registry

The container registry to use when pulling images.

- required: `false`
- default: `docker.io`

### image

Image to use for you deployment.

- required: `true`

### ports

Ports exposed from the container to the host machine.

> 80:9000 - The left port (80) refers to the host port, the right port (9000) refers to the container port.

- required: `false`

```json
{
  "ports": {
    "80": "9000"
  }
}
```

You can optionally leave the host port **blank** and Krane will find a free port and assign it. This is especially handy to avoid **port conflicts** when scaling out a deployment.

For example to load-balance a deployment with multiple instances on a specific port

```json
{
  "scale": 3,
  "ports": {
    "": "9000"
  }
}
```

In the above configuration you'll have 3 instances of your deployment load-balanced on port **9000**. See [scale](deployment-configuration?id=scale) for more details on load-balancing.

### env

The environment variables passed to the containers part of a deployment.

> ⚠️ Environment variables should not contain any sensitive data, use [`secrets`](deployment-configuration?id=secrets) instead.

- required: `false`

```json
{
  "env": {
    "NODE_ENV": "dev",
    "PORT": "8080"
  }
}
```

### secrets

Secrets are used when you want to pass sensitive information to your deployments.

> You can add deployment secrets using the Krane [CLI](cli?id=secrets)

- required: `false`

```json
{
  "secrets": {
    "SECRET_TOKEN": "@MY_SECRET_TOKEN",
    "PROXY_API_URL": "@SOME_PROXY_API_URL"
  }
}
```

### tag

The tag used when pulling the image.

- required: `false`
- default: `latest`

### volumes

The volumes to mount from the container to the host.

- required: `false`

```json
{
  "volumes": {
    "/host/path": "/container/path"
  }
}
```

### alias

Entry alias for your deployment.

> ⚠️ Aliases require an [A Record](https://www.digitalocean.com/docs/networking/dns/how-to/manage-records/#a-records) to be created in order for redirects to work.

required: `false`

```json
{
  "alias": [
    "my-app.example.com",
    "my-app-dev.example.com",
    "my-app-mybranch.example.com"
  ]
}
```

### command

Custom command to start the containers.

- required: `false`

```json
{
  "command": "npm run start --prod"
}
```

### secured

Enable HTTPS/TLS communication to your deployment. Certificates are auto-generated via [Let's Encrypt](https://letsencrypt.org/).

- required: `false`
- default: `false`

```json
{
  "secured": true
}
```

### scale

Number of containers created for a deployment. Instances are load-balanced in a [round-robin](https://en.wikipedia.org/wiki/Round-robin_DNS) fashion.

> Tip: Setting scale to 0 removes all containers for a deployment without deleting the deployment.

- required: `false`
- default: `1`

```json
{
  "scale": 3
}
```
