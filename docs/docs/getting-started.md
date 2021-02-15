# Getting Started

![Install Krane](../assets/1-install-krane.png)

You can install Krane using this interactive script. 

It is by far the *easiest* and *fastest* way to **create** or **update** a Krane instance.

```
bash <(wget -qO- get.krane.sh)
```

Other [installation](docs/installation) methods and configurations.

![Download CLI](../assets/2-download-cli.png)

Download the Krane [CLI](docs/cli) to execute commands on a Krane instance.

```
npm i -g krane
```

Full list of [commands](docs/cli?id=commands).

![Setup Authentication](../assets/3-setup-authentication.png)

Create public and private keys used for authentication.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM' -f $HOME/.ssh/krane
```

The private key stays on the user's machine, the public key is appended to `~/.ssh/authorized_keys` where Krane is running.

![Authenticate](../assets/4-authentication.png)

When logging in, you'll be prompted for the endpoint where Krane is running and the public key you created in step 3. Once authenticated you'll be able to execute commands on that Krane instance.

To switch between Krane instances you'll have to login again.

```
krane login
```

![Deploy](../assets/5-deploy.png)

Create a file and copy the following deployment configuration

`deployment.json`

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

For a full list of configuration properties, checkout the [deployment configuration](docs/deployment) section.
