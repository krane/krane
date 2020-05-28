## krane rest api

The krane server exposes an api you can use to make api calls

You can also use [krane-cli](https://github.com/biensupernice/krane-cli) to communicate to a `krane-server`

## Endpoints

| Methods     | Path                            | Auth |
| ----------- | ------------------------------- | ---- |
| POST        | /auth                           | No   |
| SSE         | /containers/:containerID/events | No   |
| GET, POST   | /deployments                    | Yes  |
| POST        | /deployments/:name/run          | Yes  |
| GET, DELETE | /deployments/:name              | Yes  |
| ws          | /deployments/:name/events       | No   |
| GET         | /sessions                       | Yes  |
| GET         | /login                          | No   |

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

`POST /deployments`

Create a new deployment

**Body**

```json
{
  "data": {
    "name": "docker-to-node",
    "config": {
      "registry": "docker.io",
      "image": "davidcasta/docker-to-node",
      "container_port": "8080",
      "host_port": "9002"
    }
  }
}
```

---

`DELETE /deployments/:name`

Delete a deployment by name. This will also remove any docker resources.

---

`POST /deployments/:name/run`

Start a deployment by name

**Query Params**

- tag : The docker image tag used for this deployment, defaults to `latest`

** Body **
Empty

**Response** 202 Accepted

This route returns immedtely since the deployment may take some time, instead an accepted response is returned if the server acknowledges the deployment request. The request starts on its own thread, to check the status see `/deployments/:name` which returns the deployment and the containers part of that deployment along with the containers status.

---

`GET /deployments/:name`

Get a deployment by name. Returns the deployment template and the containers part of that deployment. The containers include the status.

**Response**

```json
{
  "data": {
    "template": {
      "name": "docker-to-node",
      "config": {
        "registry": "docker.io",
        "image": "davidcasta/docker-to-node",
        "container_port": "8080",
        "host_port": "9002"
      }
    },
    "containers": [
      {
        "Id": "3891e40a6c97aa8bb0ebfdf13045517f7035a26254fbe5c532f9fff9dd8a2f72",
        "Names": ["/docker-to-node-9ca24de8"],
        "Image": "docker.io/davidcasta/docker-to-node:latest",
        "ImageID": "sha256:7177cc313686dad5edc09276ef4c86a3eba0e96bc8144bc3ebfba8f6ca58e7d4",
        "Command": "npm run start",
        "Created": 1590447256,
        "Ports": [
          {
            "IP": "0.0.0.0",
            "PrivatePort": 8080,
            "PublicPort": 9002,
            "Type": "tcp"
          }
        ],
        "Labels": {
          "deployment.name": "docker-to-node"
        },
        "State": "running",
        "Status": "Up 39 minutes",
        "HostConfig": {
          "NetworkMode": "default"
        },
        "NetworkSettings": {
          "Networks": {
            "krane": {
              "IPAMConfig": null,
              "Links": null,
              "Aliases": null,
              "NetworkID": "0906655e3c38fad929bde35c6e23495c2c4436eb73375777a9d3da67fd7101f4",
              "EndpointID": "9034b03d379d030692980355b78e50f9c88ed1b737a7bb4d7c0344c637078bf0",
              "Gateway": "172.24.0.1",
              "IPAddress": "172.24.0.2",
              "IPPrefixLen": 16,
              "IPv6Gateway": "",
              "GlobalIPv6Address": "",
              "GlobalIPv6PrefixLen": 0,
              "MacAddress": "02:42:ac:18:00:02"
            }
          }
        },
        "Mounts": []
      }
    ]
  }
}
```

---

`ws /deployments/:name/events`

This connection returns live events for a deployment with the following structure

Example to listen to events: [gist](https://gist.github.com/david-castaneda/b5b2f05d3ea1080692f221fb423cd344)

**Event**

```ts
type Event {
  Timestamp string
  Message string
  Deployment Template
}
```

```ts
type Template {
  Name string
  Config TemplateConfig
}
```

```ts
type TemplateConfig {
  Registry string
  Image string
  ContainerPort string
  HostPort string
}
```

---

`SSE /container/:containeID/events`

This route emits [server-sent-events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) for a container.

Example html gist: [here](https://gist.github.com/david-castaneda/c6dd5f23c8c7fbdf792b478be78fcdd4) for listening to container events

---

`GET /sessions`

List all the sessions for the server. A session is created when you log in and authenticate with the krane-server.

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

---
