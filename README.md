## krane
![Go](https://github.com/biensupernice/krane-server/workflows/Go/badge.svg?branch=master)

Easy container deployments

## Commands

| Command      | Description                    |
| ------------ | ------------------------------ |
| krane deploy | Deploy your app                |
| krane login  | Authenticate with krane server |

krane, inspired by [now](https://vercel.com/), [exoframe](https://github.com/exoframejs/exoframe), [render](https://render.com/), [dokku](http://dokku.viewdocs.io/dokku/)... is a tool for easily deploying docker apps to the cloud.

## Getting Started

You'll need the [krane-server](https://github.com/biensupernice/krane-server) installed on an inexpensive server.

```bash
curl -sf https://raw.githubusercontent.com/biensupernice/krane-server/master/bootstrap.sh | sh
```

Deploy your project using the [cli](https://github.com/biensupernice/krane-cli)

```shell
npx krane-cli deploy
```
