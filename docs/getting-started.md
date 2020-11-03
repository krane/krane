# Introduction

Krane is a container deployment tool that makes it easy to create and manage small development workloads. Krane sits on any server and interfaces with Docker exposing a productive toolset for managing containers. The Krane CLI allows you to create, configure, and automate application resources from any machine - Actions, CI, localhost.

![Krane](https://user-images.githubusercontent.com/21694364/89133914-371a5900-d4ee-11ea-9e7d-3ff5282c30f5.png)

<!-- tabs:start -->

# ** Install Krane **

Install Krane on any machine that has Docker.

```
docker run -d --name=krane \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -p 8500:8500 biensupernice/krane
```

Other [installation](installation) methods and configurations.

# ** Download CLI **

Download the Krane [CLI](cli) to communicate with Krane.

```
npm i -g @krane/cli
```

Full list of [commands](cli?id=commands).

# ** Setup Authentication **

Create public and private keys used for authentication.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'
```

The private key is kept on the user's machine, the public key is stored where Krane is running and appended to `~/.ssh/authorized_keys`

**Authenticate**

```
krane login
```

#### ** Deploy **

Create a deployment configuration file `krane.json`, example hello world deployment.

```
{
  "name": "hello-world-app",
  "image": "hello-world",
  "alias": ["hello.localhost"]
}
```

```
krane deploy -f /path/to/krane.json
```

Awe yea, first deployment 🥳

<!-- tabs:end -->
