# Krane server installation

## Installing using docker

```bash
docker run --rm --name=krane \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -v ~/.krane:/root/.krane \
    -p 80:8080 krane --build
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
