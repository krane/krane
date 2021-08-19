<img src="docs/assets/krane-wordmark.png" width="350">

> Open-source, self-hosted, container management solution

[![CI](https://github.com/krane/krane/workflows/CI/badge.svg?branch=main)](https://github.com/krane/krane/actions)
[![Release](https://img.shields.io/github/v/release/krane/krane)](https://github.com/krane/krane/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/krane/krane)](https://goreportcard.com/report/github.com/krane/krane)

Krane is a container management solution that helps you to deploy infrastructure with ease. Lightweight and easy to setup, Krane is great for developers who want to self-host infrastructure at the lowest cost possible.

- **Documentation:** https://docs.krane.sh
- **Releases:** https://github.com/krane/krane/releases
- **Bugs:** https://github.com/krane/krane/issues

## Tooling

These development tools help manage and automate infrastructure running on Krane.

- **Deployment CLI:** https://github.com/krane/cli
- **Deployment Status Page:** https://github.com/krane/statuspage
- **GitHub Action:** https://github.com/krane/action

## Features

- Krane runs on compute as low as $3.50
- Single command deployments
- Single file deployments
- Deployment DNS [aliases](https://docs.krane.sh/#/docs/deployment?id=alias) (`subdomain.example.com`)
- Deployment [secrets](https://docs.krane.sh/#/docs/deployment?id=secrets) for hiding sensitive environment variables
- Deployment [scaling](https://docs.krane.sh/#/docs/deployment?id=scale) to distribute the workload between containers
- Deployment [rate limit](https://docs.krane.sh/#/docs/deployment?id=rate_limit) to limit incoming requests
- HTTPS/TLS out-of-the-box with auto generated [Let's Encrypt](https://letsencrypt.org/) certificates
- [Self-hosted](#motivation) - Cost-effective, bring your own server, scale when you need

## Quick-Start - Install Script

1. Install Krane

The `install.sh` script provides a convenient way to download Krane on virtually any compute such as Vultr, Digital Ocean, AWS, Azure, GCP, Linode, and even on your localhost.

To install Krane just run:

```
bash <(wget -qO- get.krane.sh)
```

2. Create a deployment configuration file

`deployment.json`

```json
{
  "name": "krane-getting-started",
  "image": "docker/getting-started",
  "alias": ["getting-started.example.com"]
}
```

3. Deploy

```
krane deploy -f ./deployment.json
```

Please see the [official docs site](https://docs.krane.sh) for complete documentation.

<a name="motivation"></a>

## Motivation

Krane is a self-hosted PaaS. You bring your own server and install Krane on it to manage your containers in the form of deployments - The benefit, <i>cost per deployment</i>. A self-hosted solution allows you to own your server (cost-effective), and the benefit of any number of deployments at no extra cost. Maintaining and managing your own solution may sound complex, Krane tries to make the process <i>straight-forward</i> and <i>cost-effective</i> .

Krane isn't a replacement for [Kubernetes](https://kubernetes.io), [ECS](https://aws.amazon.com/ecs/), or any other container orchestration solution you might see running production applications, instead it's a tool you can leverage to make development of side-projects or small workloads cheap and straight forward. That was the main objective, a productive deployment tool for managing non-critical container workloads on remote servers.

## Contributions

Krane is released under the [MIT license](https://github.com/krane/krane/blob/refactor-readme/LICENSE). Please refer to the [contribution guidelines](https://github.com/krane/krane/blob/refactor-readme/CONTRIBUTING.md) before raising an issue or feature request. We appreciate all contributions, small or large, and look forward to hearing feedback and improvement proposals.
