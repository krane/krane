<img src="docs/assets/krane-wordmark.png" width="350">

> Open-source, self-hosted, container management solution

[![CI](https://github.com/krane/krane/workflows/CI/badge.svg?branch=master)](https://github.com/krane/krane/actions)
[![Release](https://img.shields.io/github/v/release/krane/krane)](https://github.com/krane/krane/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/krane/krane)](https://goreportcard.com/report/github.com/krane/krane)

Krane is a container management solution that abstracts away the hard parts from deploying infrastructrue at the lowest cost possible.

- **Documentation:** https://krane.sh
- **Releases:** https://github.com/krane/krane/releases
- **Bugs:** https://github.com/krane/krane/issues

## Tooling

- **Deployment CLI:** https://github.com/krane/cli
- **Deployment UI:** https://github.com/krane/ui
- **GitHub Action:** https://github.com/krane/action

## Features

- Single command deployments
- Single file deployments
- Deployment DNS [aliases](https://www.krane.sh/#/docs/deployment?id=alias) (`subdomain.example.com`)
- Deployment [secrets](https://www.krane.sh/#/docs/deployment?id=secrets) for hiding sensitive environment variables
- Deployment [scaling](https://www.krane.sh/#/docs/deployment?id=scale) to distribute the workload between containers
- Deployment [rate limit](https://www.krane.sh/#/docs/deployment?id=rate_limit) to limit incoming requests
- HTTPS/TLS support with auto generated [Let's Encrypt](https://letsencrypt.org/) certificates
- [Self-hosted](#motivation) - Cost-effective, bring your own server, scale when you need

## Getting started

1. Install Krane

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

[Additional deployment configuration options](https://www.krane.sh/#/docs/deployment)

<a name="motivation"></a>

## Motivation

Krane is a self-hosted PaaS. You bring your own server and install Krane on it to manage your containers in the form of deployments - The benefit, <i>cost per deployment</i>. A self-hosted solution allows you to own your server (cost-effective), and the benefit of any number of deployments at no extra cost. Maintaining and managing your own solution may sound complex, Krane tries to make the process <i>straight-forward</i> and <i>cost-effective</i> .

Krane isn't a replacement for [Kubernetes](https://kubernetes.io), [ECS](https://aws.amazon.com/ecs/), or any other container orchestration solution you might see running production applications, instead it's a tool you can leverage to make development of side-projects or small workloads cheap and straight forward. That was the main objective, a productive deployment tool for managing non-critical container workloads on remote servers.

## Contributions

Krane is released under the [MIT license](https://github.com/krane/krane/blob/refactor-readme/LICENSE). Please refer to [contribution guidelines](https://github.com/krane/krane/blob/refactor-readme/CONTRIBUTING.md) before raising an issue or feature request, we appreciate all contributions, small or large, and look forward to hearing feeback and improvement proposals.
