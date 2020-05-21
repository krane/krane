# Krane server installation

## Getting Started

You'll need the [krane-server](https://github.com/biensupernice/krane-server) installed on an inexpensive server.

```bash
curl -sf https://raw.githubusercontent.com/biensupernice/krane-server/master/bootstrap.sh | sh
```

Deploy your project using the [cli](https://github.com/biensupernice/krane-cli)

```shell
npx krane-cli deploy
```

## Commands

| Command      | Description                    |
| ------------ | ------------------------------ |
| krane deploy | Deploy an app                  |
| krane login  | Authenticate with krane server |

## Runing with docker

```bash
docker run --name=krane \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -v ~/.krane:/root/.krane \
    -p 80:8080 biensupernice/krane
```

## Creating authentication keys

```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com" -m 'PEM'

-t type
-b bytes
-C comments
-m key format
```

Now grab the contents of `key.pub` and add it to the `authorized_keys` on your server
