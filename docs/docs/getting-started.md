# Getting Started

Krane makes it easy to deploy containers for development workloads on remote or local servers. Krane interfaces with Docker exposing a productive toolset for managing containerized services known as deployments. The Krane [CLI](https://www.krane.sh/#/docs/cli) allows you to interact with Krane to create manage and automate deployments.

![Krane](https://res.cloudinary.com/biensupernice/image/upload/v1609389359/architecture_img_whesih.png)

![Install Krane](../assets/1-install-krane.png)

Install and run Krane using Docker.

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY=changeme \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh:/root/.ssh  \
    -p 8500:8500 biensupernice/krane
```

Other [installation](docs/installation) methods and configurations.

![Download CLI](../assets/2-download-cli.png)

Download the Krane [CLI](docs/cli) to execute commands on a Krane instance.

```
npm i -g @krane/cli
```

Full list of [commands](docs/cli?id=commands).

![Setup Authentication](../assets/3-setup-authentication.png)

Create public and private keys used for authentication.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'
```

The private key stays on the user's machine, the public key is appended to `~/.ssh/authorized_keys` where Krane is running.

![Authenticate](../assets/4-authentication.png)

When logging in, you'll be prompted for the endpoint where Krane is running and the public key you created in step 3. Once authenticated you'll be able to execute commands on that Krane instance.

To switch between Krane instances you'll have to login again.

```
krane login
```

![Deploy](../assets/5-deploy.png)

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

For a full list of configuration options, checkout the [deployment configuration](docs/deployment-configuration) section.
