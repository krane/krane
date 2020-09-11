# Installation

Mac

```bash
curl -L https://github.com/biensupernice/krane/releases/download/{version}/krane_{version}_darwin_amd64.tar.gz | tar xz && chmod +x krane
```

Linux

```bash
curl -L https://github.com/biensupernice/krane/releases/download/{version}/krane_{version}_linux_386.tar.gz | tar xz && chmod +x krane
```

Docker

```bash
docker run -d --name=krane \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -p 8500:8500 biensupernice/krane
```

Next step: [Authentication](authentication.md)
