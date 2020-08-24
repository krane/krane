# Krane - ![Test master](https://github.com/biensupernice/krane/workflows/test/badge.svg?branch=master)

<p align="center">
    <a href="https://github.com/biensupernice/krane">
        <img align="center" src="https://user-images.githubusercontent.com/21694364/89133914-371a5900-d4ee-11ea-9e7d-3ff5282c30f5.png" width="700"/>
    </a>
</p>

## Overview

---

> ⚠️ The development of krane is still a work in progess (WIP)

Krane is a self-hosted container management solution that runs on your hardware, whether its a linux server on any cloud provider or localhost, to interface with the Docker Engine and expose a simple API that the krane cli uses to manage your containers. The <a href="https://github.com/biensupernice/krane-cli">krane-cli</a> allows you to authenticate with krane to create container resources on the host machine.

## Installing

### Mac
```sh
curl -L https://github.com/biensupernice/krane/releases/download/{version}/krane_{version}_darwin_amd64.tar.gz | tar xz && chmod +x krane
```

### Linux
```sh
curl -L https://github.com/biensupernice/krane/releases/download/{version}/krane_{version}_linux_386.tar.gz| tar xz && chmod +x krane
```

### Docker
```sh
docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8500:8500 biensupernice/krane
```
