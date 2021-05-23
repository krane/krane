A collection of commonly deployed services using Krane, and their [deployment configurations](http://krane.sh/#/docs/deployment)

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

#### Postgres

https://www.postgresql.org

```json
{
  "name": "postgres",
  "image": "library/postgres",
  "alias": ["postgres.example.com"],
  "secure": true,
  "env": {
    "POSTGRES_DB": "pg",
    "POSTGRES_PASSWORD": "pg",
    "POSTGRES_USER": "pg"
  },
  "ports": {
    "5432": "5432"
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

#### Kafka

https://kafka.apache.org

```json
{
  "name": "kafka",
  "image": "confluentinc/cp-kafka",
  "targetPort": "29092",
  "ports": {
    "29092": "29092"
  },
  "env": {
    "KAFKA_BROKER_ID": "1",
    "KAFKA_ZOOKEEPER_CONNECT": "<ZOOKEEPER_CONTAINER_ID>:2181",
    "KAFKA_ADVERTISED_LISTENERS": "PLAINTEXT://localhost:9092,PLAINTEXT_HOST://localhost:29092",
    "KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT",
    "KAFKA_INTER_BROKER_LISTENER_NAME": "PLAINTEXT",
    "KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1"
  }
}
```

#### Zookeeper

https://zookeeper.apache.org

```json
{
  "name": "zookeeper",
  "image": "confluentinc/cp-zookeeper",
  "ports": {
    "2181": "2181"
  },
  "env": {
    "ZOOKEEPER_CLIENT_PORT": "2181",
    "ZOOKEEPER_TICK_TIME": "2000"
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

https://www.portainer.io

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

#### Supabase

https://supabase.io

```json
{
  "name": "supabase",
  "image": "supabase/postgres",
  "alias": ["supabase.example.com"],
  "secure": true,
  "secrets": {
    "POSTGRES_DB": "@POSTGRES_DB",
    "POSTGRES_PASSWORD": "@POSTGRES_PASSWORD",
    "POSTGRES_USER": "@POSTGRES_USER"
  },
  "ports": {
    "5432": "5432"
  }
}
```
