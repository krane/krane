# Status Page

The [Krane status sage](https://github.com/krane/ui) is an observability dashboard exposing the status of your Krane deployments.

<span class="img-wrapper">![Status Page](../assets/ui-page.png)</span>

## Deploying

The status page is packaged up into a [Docker image](https://hub.docker.com/repository/docker/biensupernice/krane-ui) you can directly deploy using Krane.

`deployment.json` 
 ```json
 {
   "name": "krane-statuspage",
   "image": "biensupernice/krane-ui",
   "secure": true,
   "alias": ["status.example.com"],
   "secrets": {
     "KRANE_ENDPOINT": "https://krane.example.com",
     "KRANE_TOKEN": "@KRANE_TOKEN"
   }
 }
```

```
krane deploy -f /path/to/deployment.json
```

> It's recommended to use [secrets](http://docs.krane.sh/#/docs/deployment?id=secrets) to protect against plain-text access tokens
    
## FAQ

##### How do i get a `KRANE_TOKEN`?

You can create ad-hoc access tokens using [Krane sessions](http://docs.krane.sh/#/docs/cli?id=sessions)

```
krane sessions create krane-statuspage
```
