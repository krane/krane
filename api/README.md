# Krane API

> ⚠️ The development of the krane API is still a work in progess (WIP)

The Krane API exposes endpoints for creating, updating, deleting Krane and Docker resources.

## Open Endpoints

Open endpoints do not require Authentication.

- Health: `GET /`
- Login: `GET /login`
- Authenticate: `POST /auth`

## Authenticated Endpoints

Authenticated endpoints require Authentication with the Krane server.

### Spec

- Create: `POST /spec`
- Update: `GET /spec/{name}`

### Deployments

- Get all: `GET /deployments`
- Get one: `GET /deployments/{name}`
- Run a deployment: `POST /deployments/{name}`
- Delete a deployment: `DELETE /deployments/{name}`
- Stop a deployment: `POST /deployment/{name}/stop`

### Alias

- Create or update an alias: `POST /alias/{deploymentName}`
- Delete an alias: `DELETE /alias/{deploymentName}`

### Activity

- Get recent activity: `GET /activity`
