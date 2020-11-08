# CLI

The Krane [CLI](https://github.com/krane/cli) allows you to interact with Krane to create container resources.

## Installing

```
npm i -g @krane/cli
```

## Authenticating

Krane uses private and public key authentication. Both keys are used for ensuring authenticity of incoming request.

To create both a public and private key, use the following command.

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'

-t type
-b bytes
-C comments
-m key format
```

This will generate 2 different keys, a `private` & `public (.pub)` key.

The `private key` is kept on the user's machine, the `public key` is stored where Krane is running and appended to `~/.ssh/authorized_keys`

Now you can try logging in. The CLI will prompt you to select the public key your just created. This will be used for authenticating with the private key from the Krane server.

```
krane login
```

## Commands

### login

Authenticate with Krane.

```
krane login
```

### config

List the deployment configuration for a deployment.

```
krane config <deployment>
```

### delete

Delete a deployment

```
krane delete <deployment>
```

### describe

Describe a deployment. This provides details on the running containers for a single deployment.

```
krane describe <deployment>
```

### deploy

Deploy or update an application.

```
krane deploy -f /path/to/krane.json
```

### list

List all deployments.

```
krane list
```

### secrets

List all deployment secrets.

```
krane secrets list <deployment>
```

Add a deployment secret.

```
krane secrets add <deployment> -k token -v my-secret-token
```

Delete a deployment secrets

```
krane secrets delete <deployment> -k token
```
