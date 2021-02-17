A collection of commonly deployed services using Krane, and their [deployment configurations](http://localhost:3000/#/docs/deployment)

If you'd like to add a configuration submit a pull request [here](https://github.com/krane/krane/tree/master/docs/docs/example-configs.md)

#### Mongo 

https://www.mongodb.com

```json
{
  "name": "mongo",
  "image": "library/mongo",
  "alias": ["mongo.example.com"],
  "secure": true,
  "env": {
    "MONGO_INITDB_DATABASE": "example"
  },
  "ports": {
    "": "27017"
  }
}
```

#### Vault
 
https://www.vaultproject.io

```json
{
  "name": "vault",
  "image": "library/vault",
  "alias": ["vault.example.com"],
  "secure": true,
  "secrets": {
    "VAULT_DEV_ROOT_TOKEN_ID": "@VAULT_DEV_ROOT_TOKEN_ID"
  }
}
```

#### Meili

https://www.meilisearch.com

```json
{
  "name": "meili",
  "image": "getmeili/meilisearch",
  "secure": true,
  "alias": ["meili.example.com"],
  "volumes": {
    "/tmp/data.ms": "/data.ms"
  }
}
```

#### Portainer

https://www.portainer.io/

```json
{
  "name": "portainer",
  "image": "portainer/portainer-ce",
  "secure": true,
  "alias": ["portainer.example.com"],
  "target_port": "9000",
  "ports": {
    "": "8000",
    "": "9000"
  },
  "volumes": {
    "/var/run/docker.sock": "/var/run/docker.sock"
  }
}
```