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

<b>Spec</b>: This represents the structure of a deployed application.

## Installing

---

| Operating System         | Download Link                                                                                                                                                                             |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Ubuntu 18.04.3 (LTS) x64 | https://github.com/biensupernice/krane/releases/download/0.0.1/krane_0.0.1_linux_386.tar.gz https://github.com/biensupernice/krane/releases/download/0.0.1/krane_0.0.1_linux_amd64.tar.gz |
| macOS Catalina           | https://github.com/biensupernice/krane/releases/download/0.0.1/krane_0.0.1_darwin_amd64.tar.gz                                                                                            |
|                          |                                                                                                                                                                                           |

Find the appropriate download link and use the below command to install the executable

```sh
curl -L <download link> | tar xz && chmod +x krane
```

Alternatively you can run krane using docker. It uses a lightweight Alpine image to reduce security risks with enough functionality for developing and debugging.

```sh
docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8080:8080 biensupernice/krane
```
