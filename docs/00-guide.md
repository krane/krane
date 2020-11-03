The project is open-source, in its early stages, and is no replacement for production ready solutions. So why even use Krane? Chances are most of the time when developing on side projects we come across a moment in which we wish to grab our code, put it in a container, and ship off to the cloud. If this is you, Krane might come in handy.

Krane is self-hosted, this means you can run it on your own hardware - localhost, digital ocean, aws, etc.. Theres obviously quite a lot to consider when managing and hosting your own Krane instance but we've tried to simplify the experience as much as we could. We simplified the language of describing containers just a bit, if your familiar with docker-compose then you are already more than qualified to use Krane.

Lets go over some of the basic:

**Deployments**

You can think of a single deployment as an instance of your application. For example, lets pretend we're making an app that recommends restaurants to users based on their location. Our app is gonna have a frontend, backend, database, and a recommendations service. We'll name our app _Locally_, in Krane we can model each component of Locally as a deployment, it might look something like

Insert image:

The diagram above would then look like this when configuring it with Krane

`frontend.krane.json`

```json
{
  "name": "locally-web",
  "image": "locally/web",
  "alias": ["locally.com"],
  "env": {
    "API_URL": "https://api.locally.com"
  }
}
```

`backend.krane.json`

```json
{
  "name": "locally-api",
  "image": "locally/api",
  "alias": ["api.locally.com"],
  "env": {
    "RECOMMENDATION_API_URL": "https://re.locally.com"
  },
  "secrets": {
    "DB_HOST": "@DB_HOST",
    "DB_USER": "@DB_USER",
    "DB_PASS": "@DB_PASSWORD"
  }
}
```

`database.krane.json`

```json
{
  "name": "locally-db",
  "image": "mongo",
  "tag": "3.6.20-xenial",
  "alias": ["db.locally.com"]
}
```

`recommendation.krane.json`

```json
{
  "name": "locally-recommendation",
  "image": "locally/re",
  "tag": "0.1.5",
  "alias": ["re.locally.com"]
}
```

In totals theres 4 services so we have 4 configs, the configurations are known as Krane configs. They are very similar to docker-compose files but are json formatted and offer simpler language when configuring thing like aliases and secrets. Krane uses this config to create a deployment, a deployment can have 1 or more containers running, for now we'll assume we are running 1 container per deployment. So how do we use these configs to tell Krane to ship our stuff? For that we have the Krane CLI.

The CLI is how you talk to Krane. Through the CLI we can tell Krane to use the config files we defined above to create our container resources.

For example if we wanted to deploy our database

```json
$ krane deploy -f database.krane.json
```

Or if we wanted to list all our deployments

```json
$ krane list
Name                     Image                         Tag
locally-web              locally/web                   latest
locally-api              locally/web                   latest
locally-recommendation   locally/re                    0.1.5
locally-db               mongo                         3.6.20-xenial
```

We can even add secrets to through the CLI

```json
$ krane secrets add locally-api -key DB_PASSWORD -value p@ssword
```

Simple, it might not cover all use cases, but the goal is to provide a simple and productive toolset.
