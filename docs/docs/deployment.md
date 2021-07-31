# Deployments

Krane deployments are defined in a **single file** which describe the configuration used to create container resources. A deployment configuration gives you a predictable and declarative way to manage container workloads.

> A common pattern is to have a `deployment.json` at the root of your project

`deployment.json`

```json
{
  "name": "krane-getting-started",    
  "image": "docker/getting-started",
  "alias": ["getting-started.example.com"]
}
```

The above **deployment configuration** sets up:

1. A deployment called **hello-world-app**
2. A container running the image **hello-world**
3. An alias **hello.example.com**

Check out these other [deployment configurations](http://krane.sh/#/docs/example-configs)

---

> Note: `name` and `image` are the only required properties

## name

The name of your deployment.

- required: `true`

## image

Container image to use for you deployment.

- required: `true`

## registry

The container registry credentials to use for pulling images.

- required: `false`

```json
{
  "registry": {
    "url": "docker.io",
    "username": "username",
    "password": "password"
  }
}
```

> ⚠️You should not be storing credentials in plain-text, use Krane [`secrets`](docs/deployment?id=secrets) instead.

Here's an example of setting registry secrets

```sh
$ krane secrets add <deployment> -k GITHUB_USERNAME -v <value>
$ krane secrets add <deployment> -k GITHUB_TOKEN -v <value>
```

And referencing them in your config

```json
{
  "registry": {
    "url": "ghcr.io",
    "username": "@GITHUB_USERNAME",
    "password": "@GITHUB_TOKEN"
  }
}
```

## tag

The tag used when pulling the image.

- required: `false`
- default: `latest`

## ports

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

In the above configuration you'll have 3 instances of your deployment load-balanced on port **9000**. See [scale](docs/deployment?id=scale) for more details on load-balancing.

## target_port

The target port to load-balance incoming traffic.

> Recommended when a deployment exposes multiple ports

- required: `false`

```json
{
  "scale": 3,
  "target_port": "8080",
  "ports": {
    "8080": "8080",
    "9200": "9200",
    "27017": "27017"
  }
}
```

In the above configuration, multiple ports are exposed from the host to the container but only port **8080** will be used to load balance incoming traffic.

## env

The environment variables passed to the containers part of a deployment.

> ⚠️ Environment variables should not contain any sensitive data, use [`secrets`](docs/deployment?id=secrets) instead.

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

Secrets are used when you want to pass sensitive information to your deployments.

> You can add deployment secrets using the Krane [CLI](docs/cli?id=secrets)

- required: `false`

```json
{
  "secrets": {
    "SECRET_TOKEN": "@MY_SECRET_TOKEN",
    "PROXY_API_URL": "@SOME_PROXY_API_URL"
  }
}
```

## volumes

The volumes to mount from the container to the host.

- required: `false`

```json
{
  "volumes": {
    "/host/path": "/container/path"
  }
}
```

## alias

Entry alias for your deployment.

> ⚠️ Aliases require an [A Record](https://www.digitalocean.com/docs/networking/dns/how-to/manage-records/#a-records) to be created in order for redirects to work.

required: `false`

```json
{
  "alias": [
    "example.com",
    "my-app.example.com",
    "my-app-dev.example.com",
    "my-app-mybranch.example.com"
  ]
}
```

The above configuration routes all aliases to the same deployment.

## command

Custom command to start the containers.

- required: `false`

```json
{
  "command": "npm run start --prod"
}
```

## secure

Enable HTTPS/TLS communication to your deployment. Certificates are auto-generated via [Let's Encrypt](https://letsencrypt.org/).

- required: `false`
- default: `false`

```json
{
  "secure": true
}
```

## scale

Number of containers created for a deployment. Instances are load-balanced in a [round-robin](https://en.wikipedia.org/wiki/Round-robin_DNS) fashion.

> Tip: Setting scale to 0 removes all containers for a deployment without deleting the deployment.

- required: `false`
- default: `1`

```json
{
  "scale": 3
}
```

## internal

Mark the deployment as internal. Internal deployments are used to differentiate Krane deployments from user deployments. An example of an internal deployment is the krane proxy.

- required: `false`
- default: `false`

```json
{
  "internal": true
}
```

## rate_limit

Rate limit allows you to control the incoming traffic to your deployments. It lets you to define the number of **requests per second** that a service can get in a specific predefined period. 

> Note: Default behavior is **no rate limit**  

- required: `false`
- default: `0`  which means no rate limit

```json
{
  "rate_limit": 100
}
```
