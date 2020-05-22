A deployment models what your application looks like.

To create a deployment you first need a template.

### Deployment Template

| Property       | Description                                | Default Value | Required |
| -------------- | ------------------------------------------ | ------------- | -------- |
| name           | Name of the deployment                     |               | Yes      |
| registry       | Docker registry                            | docker.io     | No       |
| image          | Docker image name                          |               | Yes      |
| tag            | Docker image tag                           | latest        | No       |
| container_port | Port to expose to the host                 |               | No       |
| host_port      | Port to map from the container to the host |               | No       |

**Example**

```json
{
  "name": "node-app",
  "config": {
    "image": "davidcasta/docker-to-node",
    "container_port": "8080",
    "host_port": "80"
  }
}
```

This template is identified by its name `node-app` and uses the `davidcasta/docker-to-node:latest` image which exposes a simple express app on port 8080. The port binding is as follows`80:8080`. The containers that this template spins up will use the following naming convention `{name}-{id}` where id is a uid given to the container.

Running a template is equivalent to:

```bash
docker run -d --name {name}-{id} -p {host_port}:{container_port} {registry}/{image}:{tag}
```

The [cli](https://github.com/biensupernice/krane-cli) walks you through generating a template. Notice how the template contains no sensitive information about your deployment. This is because a template should be checked into version control and can be picked up during ci for deploying containers in ci.

### Creating a deployment

Once you have a template, you can create a deployment using the krane [cli](https://github.com/biensupernice/krane-cli).
