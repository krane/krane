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

---

| Operating System         |                                                                                                                                                                              |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Ubuntu 18.04.3 (LTS) x64 | linux_386.tar.gz or linux_amd64.tar.gz |
| macOS Catalina           | darwin_amd64.tar.gz                                                                                            |

Once you have the appropriate download link you can run the below command to install the Krane executable

```sh
curl -L <download link> | tar xz && chmod +x krane
```

Alternatively you can run krane using docker. It uses a lightweight Alpine image to reduce security risks with enough functionality for developing and debugging.

```sh
docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8500:8500 biensupernice/krane
```
