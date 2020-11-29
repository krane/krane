# Getting Started

Krane is a container deployment tool that makes it easy to create and manage small development workloads. Krane sits on any server and interfaces with Docker exposing a productive toolset for managing containers. The Krane CLI allows you to create, configure, and automate application resources from any machine - Actions, CI, localhost.

![Krane](https://user-images.githubusercontent.com/21694364/89133914-371a5900-d4ee-11ea-9e7d-3ff5282c30f5.png)

[![Install Krane](./assets/1-install-krane.png)](https://www.krane.sh/#/installation)

Install and run Krane using Docker.

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY='changeme' \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -p 8500:8500 biensupernice/krane
```

Other [installation](installation) methods and configurations.

[![Download CLI](./assets/2-download-cli.png)](https://www.krane.sh/#/cli)

Download the Krane [CLI](cli) to execute commands on a Krane instance.

```
npm i -g @krane/cli
```

Full list of [commands](cli?id=commands).

![Setup Authentication](./assets/3-setup-authentication.png)

Create public and private keys used for authentication.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'
```

The private key is kept on the user's machine, the public key is stored where Krane is running and appended to `~/.ssh/authorized_keys`

[![Authenticate](./assets/4-authentication.png)](https://www.krane.sh/#/cli?id=authenticating)

When logging in, you'll be prompted for the endpoint where Krane is running and the public key you created in step 3. Once authenticated successfully you'll be able to execute any command on that Krane instance. 

To switch between Krane instances you'll have to login again.

```
krane login
```

[![Deploy](./assets/5-deploy.png)](https://www.krane.sh/#/cli?id=deploy)

Create a deployment configuration file `deployment.json`

For example

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