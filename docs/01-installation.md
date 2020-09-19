# Installation

Linux

```
curl -L linux.krane.sh | tar xz && chmod +x krane
```

Mac

```
curl -L mac.krane.sh | tar xz && chmod +x krane
```

Docker

```
docker run -d --name=krane \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -p 8500:8500 biensupernice/krane
```
