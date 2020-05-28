A deployment models what your application looks like.

To create a deployment you first need a spec.

### Deployment spec properties

| Property       | Description                                | Default Value | Required |
| -------------- | ------------------------------------------ | ------------- | -------- |
| name           | Name of the deployment                     |               | Yes      |
| registry       | Docker registry                            | docker.io     | No       |
| image          | Docker image name                          |               | Yes      |
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

This spec is identified by its name `node-app` and uses the `davidcasta/docker-to-node` image which exposes a simple express app on port 8080. The port binding is as follows`80:8080`. The containers that this template spins up will use the following naming convention `{name}-{id}` where id is a uid given to the container.

Running a spec is equivalent to:

```bash
docker run -d --name {name}-{id} -p {host_port}:{container_port} {registry}/{image}:latest
```

The [cli](https://github.com/biensupernice/krane-cli) walks you through generating a spec. Notice how the spec contains no sensitive information about your deployment. This is because a spec should be checked into version control and can be picked up during ci for deploying containers in ci. A spec should also not change and this is why the image `tag` was not included, rather it uses `latest` unless provided when running a deployment.

### Creating a deployment

Once you have a spec, you can create a deployment using the krane [cli](https://github.com/biensupernice/krane-cli).

### Running a deployment

Once you have a spec created on the krane-server, you can run this spec
