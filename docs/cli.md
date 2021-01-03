# CLI

The Krane [CLI](https://github.com/krane/cli) allows you to interact with Krane to create container resources.

## Installing

```
npm i -g @krane/cli
```

## Authenticating

Krane uses [private and public key authentication](https://en.wikipedia.org/wiki/Public-key_cryptography). Both keys are used for ensuring authenticity of incoming request.

1. Create the public and private key

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -f $HOME/.ssh/krane -m 'PEM'

-t type
-b bytes
-C comments
-m key format
-f output file
```

This will generate 2 different keys, a `private` & `public (.pub)` key.

2. Place the public key on the host machine where Krane is running appended to `~/.ssh/authorized_keys`.

The `private key` is kept on the user's machine.

Now you can try authenticating. The CLI will prompt you to select the public key you just created. This will be used for authenticating with the private key located on the Krane server.

```
krane login
```

## Commands

### login

Authenticate with a Krane instance.

```
krane login
```

### delete

Delete a deployment.

```
krane delete <deployment>
```

### describe

Describe a deployment in detail. This provides details on the containers part of the deployment.

```
krane describe <deployment>
```

### deploy

Create or run a deployment.

```
krane deploy -f </path/to/deployment.json>
```

Flags:

- `--file`(`-f`): Path to deployment configuration

- `--tag`(`-t`): Image tag to apply to the deployment

- `--scale`(`-s`): Number of containers to create (`default` is 1 container)

### list

List all deployments.

```
krane list
```

### logs

Stream container logs for a deployment.

```
krane logs <container>
```

### history

Get recent activity for a deployment.

```
krane history <deployment>
```

### start

Start all containers for a deployment. This command will not *create* containers, only start any stopped containers.

```
krane start <deployment>
```

### stop

Stop all containers for a deployment.

```
krane stop <deployment>
```

### restart

Restart a deployment *re-creating* container resources.

```
krane restart <deployment>
```

### secrets

List all deployment secrets.

```
krane secrets list <deployment>
```

Add a deployment secret.

```
krane secrets add <deployment> -k <key> -v <value>
```

Delete a deployment secret.

```
krane secrets delete <deployment> -k <key>
```