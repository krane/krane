# Installation


Mac
```sh
curl -L https://github.com/biensupernice/krane/releases/download/{version}/krane_{version}_darwin_amd64.tar.gz | tar xz && chmod +x krane
```

Linux
```sh
curl -L https://github.com/biensupernice/krane/releases/download/{version}/krane_{version}_linux_386.tar.gz| tar xz && chmod +x krane
```
Docker
```sh
docker run -d --name=krane \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh/authorized_keys:/root/.ssh/authorized_keys  \
    -p 8500:8500 biensupernice/krane
```

Once you have installed Krane using one of the above methods, you'll need to create a pair of keys that are used for authenticating the users that will be interacting with Krane.

See: https://github.com/biensupernice/krane/wiki/Authentication

