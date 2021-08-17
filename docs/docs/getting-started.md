# Getting Started

![Install Krane](../assets/1-install-krane.png)

Run the below command to install Krane. 

```
bash <(wget -qO- get.krane.sh)
```

Other [installation](docs/installation) methods and configurations.

![Download CLI](../assets/2-download-cli.png)

Download the [Krane CLI](docs/cli) to execute commands on the Krane instance created above.

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

When login in you'll be prompted for the private key created in step 3. Once authenticated you'll be able to execute commands on that Krane instance.

```
krane login https://krane.example.com
```

![Deploy](../assets/5-deploy.png)

Create a file and copy the following deployment configuration

`deployment.json`

```json
{
  "name": "krane-getting-started",    
  "image": "docker/getting-started",
  "alias": ["getting-started.example.com"]
}
```

```
krane deploy -f /path/to/deployment.json
```

For the full list of configuration options checkout out the [deployments pages](docs/deployment).
