# Installation

Krane should be installed on a machine running Docker

> ⚠️ If you're not sure if Docker is running, run `docker -v` and verify the output (ex: `Docker version 19.03.13, build 4484c46d9d`)

## Docker

Run Krane using Docker

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY='changeme' \
    -v /path/to/authorized_keys:/root/.ssh/authorized_keys  \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -p 8500:8500 biensupernice/krane
```

## Linux

Run Krane using the executable for Linux

```
# set Krane private key
export KRANE_PRIVATE_KEY=changeme

# install the executable
curl -L linux.krane.sh | tar xz && chmod +x krane

# run the executable
krane &
```

## Mac

Run Krane using the executable for Mac

```
# set Krane private key
export KRANE_PRIVATE_KEY=changeme

# install the executable
curl -L mac.krane.sh | tar xz && chmod +x krane

# run the executable
krane &
```

---

#### Environment Configuration

The following properties can be set as environment variables when running Krane.

| Env                        | Description                                                        | Required | Default        |
| -------------------------- | ------------------------------------------------------------------ | -------- | -------------- |
| KRANE_PRIVATE_KEY          | The private key used by Krane for signing authentication requests. | true     |                |
| LISTEN_ADDRESS             | Address and port Krane will listen on                              | false    | 127.0.0.1:8500 |
| LOG_LEVEL                  | Can only be debug\|info\|warn\|error                               | false    | info           |
| DB_PATH                    | Path to boltdb                                                     | false    | /tmp/krane.db  |
| DOCKER_BASIC_AUTH_USERNAME | Username used when authenticating with Docker                      | false    |                |
| DOCKER_BASIC_AUTH_PASSWORD | Password used when authenticating with Docker                      | false    |                |