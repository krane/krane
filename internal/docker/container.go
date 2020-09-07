package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// KraneContainerLabelName : used to identify krane managed containers
const KraneContainerLabel = "krane.deployment.name"

// CreateContainerConfig : properties required to create a container
type CreateContainerConfig struct {
	Name          string
	Image         string
	NetworkID     string
	HostPort      string
	ContainerPort string
	Labels        map[string]string
	Env           map[string]string // Comma separate string env. ex: "NODE_ENV=dev"
	Volumes       map[string]string // ex: /var/run/docker.sock:/var/run/docker.sock
}

// CreateContainer : create docker container
func (c *DockerClient) CreateContainer(
	ctx *context.Context,
	conf *CreateContainerConfig,
) (container.ContainerCreateCreatedBody, error) {

	// Configure Host Port
	hostBinding := nat.PortBinding{
		// HostIP:   "localhost",
		HostPort: conf.HostPort,
	}

	// Configure Container Port
	containerPort, err := nat.NewPort("tcp", conf.ContainerPort)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	// Bind host-to-container ports
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	// Setup volumes
	volumes := make([]mount.Mount, 0)
	for s, t := range conf.Volumes {
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: s,
			Target: t,
		})
	}
	hostConf := &container.HostConfig{
		PortBindings: portBinding,
		AutoRemove:   false,
		Mounts:       volumes,
	}

	// Normalize Env vars to be represented as an array of strings &  not a map
	envars := make([]string, 0)
	for k, v := range conf.Env {
		envar := fmt.Sprintf("%s=%s", k, v) // ex. NODE_ENV=dev
		envars = append(envars, envar)
	}

	// Setup container conf
	containerConf := &container.Config{
		Hostname: conf.Name,
		Image:    conf.Image,
		Env:      envars,
		Labels:   conf.Labels,
	}

	// Setup networking conf
	endpointConf := map[string]*network.EndpointSettings{"krane": &network.EndpointSettings{NetworkID: conf.NetworkID}}
	networkConf := &network.NetworkingConfig{EndpointsConfig: endpointConf}

	return c.ContainerCreate(*ctx, containerConf, hostConf, networkConf, conf.Name)
}

// StopContainer : stop docker container
func (c *DockerClient) StopContainer(ctx *context.Context, containerID string) error {
	timeout := 60 * time.Second
	return c.ContainerStop(*ctx, containerID, &timeout)
}

// RemoveContainer : remove docker container
func (c *DockerClient) RemoveContainer(ctx *context.Context, containerID string) error {
	options := types.ContainerRemoveOptions{}
	return c.ContainerRemove(*ctx, containerID, options)
}

func (c *DockerClient) GetOneContainer(ctx *context.Context, containerId string) (types.ContainerJSON, error) {
	return c.ContainerInspect(*ctx, containerId)
}

// GetAllContainers : gets all containers on the host machine
func (c *DockerClient) GetAllContainers(ctx *context.Context) (containers []types.Container, err error) {
	options := types.ContainerListOptions{
		All:   true,
		Quiet: false,
	}

	return c.ContainerList(*ctx, options)
}

// GetContainerStatus : get the status of a container
func (c *DockerClient) GetContainerStatus(ctx *context.Context, containerID string, stream bool) (stats types.ContainerStats, err error) {
	return c.ContainerStats(*ctx, containerID, stream)
}

func GetContainers(ctx *context.Context, deploymentName string) ([]types.Container, error) {
	client := GetClient()

	// Find all containers
	allContainers, err := client.GetAllContainers(ctx)
	if err != nil {
		return make([]types.Container, 0), err
	}

	deploymentContainers := make([]types.Container, 0)
	for _, container := range allContainers {
		kraneLabel := container.Labels[KraneContainerLabel]
		if kraneLabel == deploymentName {
			deploymentContainers = append(deploymentContainers, container)
		}
	}

	return deploymentContainers, nil
}

// ReadContainerLogs :
func (c *DockerClient) ReadContainerLogs(ctx *context.Context, containerID string) (reader io.Reader, err error) {
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
func (c *DockerClient) StartContainer(ctx *context.Context, containerID string) (err error) {
	options := types.ContainerStartOptions{}
	return c.ContainerStart(*ctx, containerID, options)
}

// ConnectContainerToNetwork : connect a container to a network
func (c *DockerClient) ConnectContainerToNetwork(ctx *context.Context, networkID string, containerID string) (err error) {
	config := network.EndpointSettings{
		NetworkID: networkID,
	}
	return c.NetworkConnect(*ctx, networkID, containerID, &config)
}
