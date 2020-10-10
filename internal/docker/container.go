package docker

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

// CreateContainerConfig : properties required to create a container
type CreateContainerConfig struct {
	ContainerName string
	Image         string
	NetworkID     string
	Labels        map[string]string
	Ports         nat.PortMap
	Volumes       []mount.Mount
	Env           []string // Comma separated, formatted NODE_ENV=dev
}

// create a docker container
func (c *Client) CreateContainer2(
	ctx context.Context,
	config CreateContainerConfig,
) (container.ContainerCreateCreatedBody, error) {
	networkingConfig := makeNetworkingConfig(config.NetworkID)
	containerConfig := makeContainerConfig(config.ContainerName, config.Image, config.Env, config.Labels)
	hostConfig := makeHostConfig(config.Ports, config.Volumes)

	return c.ContainerCreate(
		ctx,
		&containerConfig,
		&hostConfig,
		&networkingConfig,
		config.ContainerName,
	)
}

func makeContainerConfig(hostname string, image string, env []string, labels map[string]string) container.Config {
	return container.Config{
		Hostname: hostname,
		Image:    image,
		Env:      env,
		Labels:   labels,
	}
}

func makeHostConfig(ports nat.PortMap, volumes []mount.Mount) container.HostConfig {
	return container.HostConfig{
		PortBindings: ports,
		AutoRemove:   true,
		Mounts:       volumes,
	}
}

// CreateContainer : create docker container
// func (c *DockerClient) CreateContainer(
// 	ctx *context.Context,
// 	conf *CreateContainerConfig,
// ) (container.ContainerCreateCreatedBody, error) {
//
// 	// Configure Host Port
// 	hostBinding := nat.PortBinding{
// 		// HostIP:   "localhost",
// 		HostPort: conf.HostPort,
// 	}
//
// 	// Configure Container Port
// 	containerPort, err := nat.NewPort(string(container.TCP), conf.ContainerPort)
// 	if err != nil {
// 		return container.ContainerCreateCreatedBody{}, err
// 	}
//
// 	// Bind host-to-container ports
// 	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
//
// 	// Setup volumes
// 	volumes := make([]mount.Mount, 0)
// 	for s, t := range conf.Volumes {
// 		volumes = append(volumes, mount.Mount{
// 			Type:   mount.TypeBind,
// 			Source: s,
// 			Target: t,
// 		})
// 	}
// 	hostConf := &container.HostConfig{
// 		PortBindings: portBinding,
// 		AutoRemove:   false,
// 		Mounts:       volumes,
// 	}
//
// 	// Normalize Env vars to be represented as an array of strings &  not a map
// 	envars := make([]string, 0)
// 	for k, v := range conf.Env {
// 		envar := fmt.Sprintf("%s=%s", k, v) // ex. NODE_ENV=dev
// 		envars = append(envars, envar)
// 	}
//
// 	// Setup container conf
// 	containerConf := &container.Config{
// 		Hostname: conf.Name,
// 		Image:    conf.Image,
// 		Env:      envars,
// 		Labels:   conf.Labels,
// 	}
//
// 	// Setup networking conf
// 	endpointConf := map[string]*network.EndpointSettings{"krane": &network.EndpointSettings{NetworkID: conf.NetworkID}}
// 	networkConf := &network.NetworkingConfig{EndpointsConfig: endpointConf}
//
// 	return c.ContainerCreate(*ctx, containerConf, hostConf, networkConf, conf.Name)
// }

// StopContainer : stop docker container
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	timeout := 60 * time.Second
	return c.ContainerStop(ctx, containerID, &timeout)
}

// RemoveContainer : remove docker container
func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	options := types.ContainerRemoveOptions{Force: force}
	return c.ContainerRemove(ctx, containerID, options)
}

func (c *Client) GetOneContainer(ctx context.Context, containerId string) (types.ContainerJSON, error) {
	return c.ContainerInspect(ctx, containerId)
}

// GetAllContainers : gets all containers on the host machine
func (c *Client) GetAllContainers(ctx *context.Context) (containers []types.Container, err error) {
	options := types.ContainerListOptions{
		All:   true,
		Quiet: false,
	}

	return c.ContainerList(*ctx, options)
}

// GetContainerStatus : get the status of a container
func (c *Client) GetContainerStatus(ctx context.Context, containerID string, stream bool) (stats types.ContainerStats, err error) {
	return c.ContainerStats(ctx, containerID, stream)
}

func (c *Client) GetContainers(ctx *context.Context, deploymentName string) ([]types.Container, error) {
	// Find all containers
	allContainers, err := c.GetAllContainers(ctx)
	if err != nil {
		return make([]types.Container, 0), err
	}

	deploymentContainers := make([]types.Container, 0)
	for _, container := range allContainers {
		kraneLabel := container.Labels["TODO"]
		if kraneLabel == deploymentName {
			deploymentContainers = append(deploymentContainers, container)
		}
	}

	return deploymentContainers, nil
}

func (c *Client) FilterContainersByDeployment(deploymentName string) ([]types.Container, error) {
	deploymentContainers := make([]types.Container, 0)

	ctx := context.Background()
	containers, err := c.GetAllContainers(&ctx)
	ctx.Done()

	if err != nil {
		logrus.Errorf("Unable to filter container by deployment, %s", err.Error())
		return make([]types.Container, 0), err
	}

	for _, container := range containers {
		kraneLabel := container.Labels["TODO"]
		if kraneLabel == deploymentName {
			deploymentContainers = append(deploymentContainers, container)
		}
	}

	return deploymentContainers, nil
}

// ReadContainerLogs :
func (c *Client) ReadContainerLogs(ctx *context.Context, containerID string) (reader io.Reader, err error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "50",
	}

	return c.ContainerLogs(*ctx, containerID, options)
}

// StartContainer : start docker container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	options := types.ContainerStartOptions{}
	return c.ContainerStart(ctx, containerID, options)
}

// ConnectContainerToNetwork : connect a container to a network
func (c *Client) ConnectContainerToNetwork(ctx *context.Context, networkID string, containerID string) (err error) {
	config := network.EndpointSettings{NetworkID: networkID}
	return c.NetworkConnect(*ctx, networkID, containerID, &config)
}
