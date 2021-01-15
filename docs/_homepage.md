<!-- ![Logo](assets/logo.png) -->

<h2><b>KRANE</b></h2> Toolset for working with containers

<span class="img-wrapper-no-shadow">![Krane](./assets/krane-flow-horizontal.png)</span>

[Krane](https://github.com/krane) makes it easy to deploy containers on remote or local servers by interfacing with Docker to expose a productive toolset for managing containerized applications in the form of deployments. The Krane [CLI](docs/cli) allows you to interact with a Krane instance to run deployments, read container logs, store deployment secrets and more. The Krane [GitHub Action](https://github.com/marketplace/actions/krane) allows you to automate deployments to continuously deliver updates when changes occur to your projects.

Check out the [getting started](docs/getting-started.md) to get up and running using Krane.

#### Motivation

Krane is a self-hosted PaaS. You bring your own server and install Krane on it to manage your containers in the form of deployments - The benefit, <i>cost per deployment</i>. A self-hosted solution allows you to own your server (cost-effective), and the benefit of any number of deployments at no extra cost. Maintaining and managing your own solution may sound complex, Krane tries to make the process <i>straight-forward</i> and <i>cost-effective</i> .

Krane isn't a replacement for [Kubernetes](https://kubernetes.io/), [ECS](https://aws.amazon.com/ecs/?whats-new-cards.sort-by=item.additionalFields.postDateTime&whats-new-cards.sort-order=desc), or any other container orchestration solution you might see running production applications, instead it's a tool you can leverage to make development of side-projects or small workloads cheap and straight forward. That was the main objective, a productive deployment tool for managing non-critical container workloads on remote servers.

