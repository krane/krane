# CLI

The Krane [CLI](https://github.com/krane/cli) allows you to interact with a Krane instance to run deployments, read container logs, store deployment secrets and more.

> Ensure you have succesfully setup your [authentication](docs/authentication.md)

## Installing

```
npm i -g krane
```

## Commands

### login

Authenticate with a Krane instance.

```
krane login
```

### list

List all deployments and their status.

```
krane list
```

### status

Returns information about the containers for a deployment.

```
krane status <deployment>
```

or get the status of a single container.

```
krane status <deployment> <container>
```

### deploy

Create or run a deployment.

```
krane deploy -f </path/to/deployment.json>
```

Flags:

- `--file`(`-f`) Path to deployment configuration

  - required: `false`
  - default: `./deployment.json`

- `--tag`(`-t`) Image tag to apply to the deployment

  - required: `false`
  - default: `latest`

- `--scale`(`-s`) Number of containers to create
  - required: `false`
  - default: `1`

### logs

Read realtime logs for a container.

```
krane logs <container>
```

### history

Get recent activity for a deployment.

```
krane history <deployment>
```

or get additional details on a specific deployment job.

```
krane history <deployment> <job>
```

### edit

Edit a deployments configuration.

```
krane edit <deployment>
```

### remove

Remove a deployment and its container resources, secrets, and configuration.

```
krane remove <deployment>
```

### start

Start all containers for a deployment. This command will not _create_ containers, only start any stopped containers.

```
krane start <deployment>
```

### stop

Stop all containers for a deployment.

```
krane stop <deployment>
```

### restart

Restart a deployment _re-creating_ container resources.

```
krane restart <deployment>
```

### secrets

List all secrets for a deployment.

```
krane secrets list <deployment>
```

Add secret to a deployment.

```
krane secrets add <deployment> -k <key> -v <value>
```

Delete a secret for a deployment.

```
krane secrets delete <deployment> -k <key>
```

### help

Display CLI usage and commands.

```
krane help
```

