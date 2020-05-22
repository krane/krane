## krane rest api

The krane server exposes an api you can use to make api calls

You can also use [krane-cli](https://github.com/biensupernice/krane-cli) to communicate to a `krane-server`

## Endpoints

| Methods   | Path               | Auth |
| --------- | ------------------ | ---- |
| POST      | /auth              | No   |
| GET       | /containers        | Yes  |
| GET, POST | /deployments       | Yes  |
| GET       | /deployments/:name | Yes  |
| GET       | /sessions          | Yes  |
| GET       | /login             | No   |

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

---

`POST /auth`

Once you have requested to login from the krane server, this route is responsible for authenticating the request by using the public key on the server to decode the signed token that contains the phrase from the server.

**Body**

```json
{
  "data": {
    "request_id": "<id from server>",
    "token": "<signed_jwt>"
  }
}
```

**Response**

```json
{
  "data": {
    "session": {
      "id": "<session_id>",
      "token": "<to authenticate subsequent requests>",
      "expires_at": "MM/DD/YYYY"
    }
  }
}
```

---

`GET /deployments`

List of all the deployments

**Response**

```json
{
  "data": [
    {
      "name": "app1",
      "config": {
        "registry": "",
        "image": "",
        "tag": "",
        "container_port": "",
        "host_port": ""
      }
    },
    {
      "name": "app2",
      "config": {
        "registry": "",
        "image": "",
        "tag": "",
        "container_port": "",
        "host_port": ""
      }
    }
  ]
}
```

---

`GET /deployments/:name`

List a single deployment by name

Assuming `:name` is app1, the response will look like the one below

**Response**

```json
{
  "data": {
    "name": "app1",
    "config": {
      "registry": "",
      "image": "",
      "tag": "",
      "container_port": "",
      "host_port": ""
    }
  }
}
```

---

`POST /deployments`

Create a new deployment

**Body**

```json
{
  "data": {
    "name": "node-app",
    "config": {
      "image": "davidcasta/docker-to-node",
      "container_port": "8080",
      "host_port": "80"
    }
  }
}
```

**Response** 202 Accepted

This route returns immedtely since the deployment may take some time, instead an accepted response is returned if the server acknowledges the deployment request. The request starts on its own thread, to check the status see `/containers`.

---

`GET /containers`

List all the containers on the server

**Response**

Example response, as you can see data contains the docker representation of the `Container` type. No reason to maintain our own container data type, this means we can leverage the data directly from the docker client on any of the krane interfaces even if the docker client used gets updates

```json
{
  "data": [
    {
      "Id": "7683abd401579f2a4f2f97b267c6fa584447383e96db8b0f49473eebaf8abbee",
      "Names": ["/app1-903c8479-05ae-5ca5-aa3d-a8e26e1ba148"],
      "Image": "docker.io/biensupernice/docker-to-node:latest",
      "ImageID": "sha256:7177cc313686dad5edc09276ef4c86a3eba0e96bc8144bc3ebfba8f6ca58e7d4",
      "Command": "npm run start",
      "Created": 1590099327,
      "Ports": [
        {
          "PrivatePort": 8080,
          "Type": "tcp"
        }
      ],
      "Labels": {},
      "State": "running",
      "Status": "Up 2 hours",
      "HostConfig": {
        "NetworkMode": "default"
      },
      "NetworkSettings": {
        "Networks": {
          "bridge": {
            "IPAMConfig": null,
            "Links": null,
            "Aliases": null,
            "NetworkID": "a30d0e134ddde5616c0b6a3316877022bbe3c622099f79e8d8f1bca9c910a08b",
            "EndpointID": "36d88ef2793cc185a0e4a2fe6aa58024cfa13960f71636cc1ebf27ddfe8f9dfc",
            "Gateway": "172.17.0.1",
            "IPAddress": "172.17.0.13",
            "IPPrefixLen": 16,
            "IPv6Gateway": "",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "MacAddress": "02:42:ac:11:00:0d"
          }
        }
      },
      "Mounts": []
    }
  ]
}
```

---

`GET /sessions`

List all the sessions for the server. A session is created when you log in.

**Response**

```json
{
  "data": [
    {
      "id": "<session_id>",
      "token": "<session_token>",
      "expires_at": "<session_expiration_date>"
    }
  ]
}
```
