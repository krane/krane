## krane rest api

The krane server exposes an api you can use to make api calls(We encourage you to use the [krane-cli](https://github.com/biensupernice/krane-cli))

## Endpoints

| Path         | Auth |
| ------------ | ---- |
| /login       | No   |
| /auth        | No   |
| /deployments | Yes  |

---

`GET /login`

This route is used when initialy logging into the krane-server. It returns a unique id and phrase. The phrase is used as the content intside the jwt token which is signed using a private key from client.

**Response**

```json
{
  "request_id": "<unique request id from server>",
  "phrase": "<unique phrase from server>"
}
```

`POST /auth`

Once you have requested to login from the krane server, this route is responsible for authenticating the request by using the public key on the server to decode the private key signed token that contains the phrase.

**Body**

```json
{
  "request_id": "<id from server>",
  "token": "<sign_jwt>"
}
```

**Response**

```json
{
  "session": {
    "id": "<session_id>",
    "token": "<for authenticating subsequent requests>",
    "expires_at": "MM/DD/YYYY"
  }
}
```
