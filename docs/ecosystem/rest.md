# REST API

- [Postman collection](postman/collection)
- [Postman collection env](postman/collection-env)

---

## Authentication

### Login

`GET /login`

The phrase is a server generated phrase that the client must sign using their private key. The signed token is used when making an `/auth` request.

**response**

```json
{
  "success": true,
  "code": 200,
  "data": {
    "request_id": "<id from server>",
    "phrase": "<server phrase>"
  }
}
```

### Authenticate

`POST /auth`

The token is a JWT token containing the server phrase signed with the clients private key.

**request body**

```json
{
  "request_id": "<id from server>",
  "token": "<signed_token>"
}
```

**response**

```json
{
  "success": true,
  "code": 200,
  "data": {
    "session": {
      "id": "<session_id>",
      "token": "<to authenticate subsequent requests>",
      "expires_at": "<mm/dd/yyyy>"
    }
  }
}
```

---

## Deployments

### Create a deployment

`POST /deployments`

**request body**

- A [Krane Config](components/krane-config)

**response**

- 200 OK

`GET /deployments`

### Get all deployments

**response**

- A list of [Krane Configs](components/krane-config)

### Get a deployment

`GET /deployments/{name}`

**response**

- A [Krane Config](components/krane-config)

`DELETE /deployment/{name}`

Delete a deployment.

**response**

- 200 OK

---

## Jobs

### Get all deployment jobs

`GET /jobs`

**response**

Array of jobs.

### Get jobs for a deployment

`GET /jobs/{name}`

**response**

Array of jobs.

### Get job by id

`GET /jobs/{name}/{id}`

**response**

A single job.
