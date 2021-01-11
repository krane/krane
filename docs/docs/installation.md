# Installation

> Note: Docker should be running on the machine where you plan on installing Krane

You can install Krane using this interactive script. 

It is by far the *easiest* and *fastest* way to **create** or **update** a Krane instance.

```
bash <(wget -qO- https://raw.githubusercontent.com/krane/krane/master/bootstrap.sh)
```

## Docker

Run Krane using Docker

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY=changeme \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh:/root/.ssh  \
    -p 8500:8500 biensupernice/krane
```

## Linux

Run Krane using the executable for Linux

```
# set Krane private key
export KRANE_PRIVATE_KEY=changeme

# install the executable
curl -L https://github.com/krane/krane/releases/download/${KRANE_VERSION}/krane_${KRANE_VERSION}_linux_386.tar.gz | tar xz && chmod +x krane

# run the executable
krane &
```

## Mac

Run Krane using the executable for Mac

```
# set Krane private key
export KRANE_PRIVATE_KEY=changeme

# install the executable
curl -L https://github.com/krane/krane/releases/download/${KRANE_VERSION}/krane_${KRANE_VERSION}_darwin_amd64.tar.gz | tar xz && chmod +x krane

# run the executable
krane &
```

---

#### Environment Configuration

The following properties can be set as environment variables when running Krane.

> Note: KRANE_PRIVATE_KEY is the only required environment variable

| Env                        | Description                                                                      | Required | Default        |
| -------------------------- | -------------------------------------------------------------------------------- | -------- | -------------- |
| KRANE_PRIVATE_KEY          | The private key used by Krane for signing authentication requests.               | true     |                |
| LISTEN_ADDRESS             | Address and port Krane will listen on                                            | false    | 127.0.0.1:8500 |
| LOG_LEVEL                  | Can only be debug\|info\|warn\|error                                             | false    | info           |
| DB_PATH                    | Path to boltdb                                                                   | false    | /tmp/krane.db  |
| DOCKER_BASIC_AUTH_USERNAME | Username used when authenticating with Docker                                    | false    |                |
| DOCKER_BASIC_AUTH_PASSWORD | Password used when authenticating with Docker                                    | false    |                |
| PROXY_ENABLED              | Enable network proxy (When disabled, aliases will not work)                      | false    | true           |
| PROXY_DASHBOARD_SECURE     | Enable HTTPS/TLS on the proxy dashboard                                          | false    | false          |
| PROXY_DASHBOARD_ALIAS      | Alias for the proxy dashboard (ex: `monitor.example.com`)                        | false    |                |
| LETSENCRYPT_EMAIL          | Email used for generating Let's Encrypt TLS certificates (must be a valid email) | false    |                |
