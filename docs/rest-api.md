## krane rest api

The krane server exposes an api you can use to make api calls

You can also use [krane-cli](https://github.com/biensupernice/krane-cli) to communicate to a `krane-server`

## Endpoints

| Path         | Auth |
| ------------ | ---- |
| /login       | No   |
| /auth        | No   |
| /deployments | Yes  |

---

`GET /login`

This route is used when initialy logging into the krane-server. It returns a unique id and phrase. The phrase is used as the data inside the jwt token which is signed using the private key choosen by the client.

**Response**

```json
{
  "request_id": "<unique request id from server>",
  "phrase": "<unique phrase from server>"
}
```

`POST /auth`

Once you have requested to login from the krane server, this route is responsible for authenticating the request by using the public key on the server to decode the signed token that contains the phrase from the server.

**Body**

```json
{
  "request_id": "<id from server>",
  "token": "<signed_jwt>"
}
```

**Response**

```json
{
  "session": {
    "id": "<session_id>",
    "token": "<to authenticate subsequent requests>",
    "expires_at": "MM/DD/YYYY"
  }
}
```

`GET /deployments`

List all the deployments and their status

**Authorization** Bearer

**Response**

```json
{
  "deployments": []
}
```

`POST /deployments`

Create a new deployment

**Body**

```json
{
  "name": "backend-service-1", // Name of the deployment
  "config": {
    "repo": "dockerhub.io",
    "image": "app/backend-service-1",
    "tag": "latest",
    "host_port": "8080",
    "container_port": "8080"
  }
}
```

**Response** 202 Accepted
