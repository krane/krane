# Configuration

Deploying containers using **Krane** starts with a configuration file that describes how Krane should run your containers.

> A recommended pattern is to have a `krane.json` file at the root of you project.

An example of a `krane.json` for deploying the Docker [hello-world](https://hub.docker.com/_/hello-world) example.

```json
{
  "name": "hello-world-app",       <- Deployment name
  "image": "hello-world",          <- Docker image
  "alias": ["hello.localhost"]     <- Custom aliases
}
```

The above deployment configuration:

- Creates a deployment called **hello-world-app**

- Uses the image **hello-world**

- Has an alias **hello.localhost** that is automatically handled by Krane.


### name

The name of your deployment.

- required: `yes`

### registry

The container registry to use when pulling images.

- required: `false`
- default: `docker.io`

### image

The image used when pulling, creating and running your deployments containers.

- required: `true`

### ports

The ports to expose from the container to the host machine.

> ⚠️ Ports are discouraged since port conflicts can become a frustating effect. Instead, Krane uses a reverse proxy that handles exposing your containers using `aliases`

- required: `false`

```json
{
  "ports": {
    "80": "8080",
  }
}
```

### env

The environment variables passed to the containers part of a deployment.

> ⚠️ Environment variables should not contain any sensitive data, use `secrets` instead.

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

Secrets are used when you want to pass sensitive information to your deployments. Secrets are **not shared** across deployments, they are only provided to the containers in the same deployment.

Secrets are created using the krane `cli` and referenced in your Krane configuration using the `@` symbol.

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


```json
{
  "alias": ["app2.example.com", "app2-dev.example.com", "app2-mybranch.example.com"]
}

```
