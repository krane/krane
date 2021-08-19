# Installation

> Note: Docker should be running on the machine where you plan on installing Krane

You can install Krane using the below command.

```
bash <(wget -qO- get.krane.sh)
```

## Docker

### Docker examples

This is the most minimal Docker example to install Krane if you preffer not to use the provided install script above.

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY=changeme \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh:/root/.ssh  \
    -p 8500:8500 ghcr.io/krane/krane
```

This is a complete Docker example with:

- Automatic HTTPS w/ Lets Encrypt certificates for deployments
- Volumed Krane DB (for storing session & deployment details)
- Log level set to debug (for debugging)

```
docker run -d --name=krane \
    -e KRANE_PRIVATE_KEY=changeme \
    -e LOG_LEVEL=debug \
    -e PROXY_ENABLED=true \
    -e PROXY_DASHBOARD_SECURE=true \
    -e PROXY_DASHBOARD_ALIAS=monitor.example.com \
    -e LETSENCRYPT_EMAIL=email@example.com \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/.ssh:/root/.ssh  \
    -v /tmp/krane.db:/tmp/krane.db \
    -p 8500:8500 ghcr.io/krane/krane
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

> Note: Krane is currently not compatible with linux/arm64/v8 machines (m1 chip)

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

| Env                        | Description                                                                                          | Required | Default        |
| -------------------------- | ---------------------------------------------------------------------------------------------------- | -------- | -------------- |
| KRANE_PRIVATE_KEY          | The private key used by Krane for signing authentication requests.                                   | true     |                |
| LISTEN_ADDRESS             | Address and port Krane will listen on                                                                | false    | 127.0.0.1:8500 |
| LOG_LEVEL                  | Can only be debug\|info\|warn\|error                                                                 | false    | info           |
| DB_PATH                    | Path to boltdb                                                                                       | false    | /tmp/krane.db  |
| PROXY_ENABLED              | Enable network proxy (When disabled, aliases will not work)                                          | false    | true           |
| PROXY_DASHBOARD_SECURE     | Enable HTTPS/TLS on the proxy dashboard                                                              | false    | false          |
| PROXY_DASHBOARD_ALIAS      | Alias for the proxy dashboard (ex: `monitor.example.com`)                                            | false    |                |
| LETSENCRYPT_EMAIL          | Email used for generating Let's Encrypt TLS certificates (must be a valid email)                     | false    |                |
| WORKERPOOL_SIZE            | Amount of workers running executing jobs. Workers run in parallel picking up jobs from the job queue | false    | 1              |
| JOB_QUEUE_SIZE             | Amount of jobs queue'd at a given time                                                               | false    | 1              |
| JOB_MAX_RETRY_POLICY       | Max retries for any job being executed                                                               | false    | 5              |
| DEPLOYMENT_RETRY_POLICY    | Max retries for a deployment                                                                         | false    | 1              |
