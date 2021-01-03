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
)

// DockerConfig properties required to create a docker container
type DockerConfig struct {
	ContainerName string
	Image         string
	NetworkID     string
	Labels        map[string]string
	Ports         nat.PortMap
	VolumeMounts  []mount.Mount
	VolumeSet     map[string]struct{}
	Env           []string // Comma separated, formatted NODE_ENV=dev
	Command       []string
	Entrypoint    []string
}

// CreateContainer creates a docker container from a Dcoker config
func (c *Client) CreateContainer(ctx context.Context, config DockerConfig) (container.ContainerCreateCreatedBody, error) {
	networkingConfig := createNetworkingConfig(config.NetworkID)
	hostConfig := createHostConfig(config.Ports, config.VolumeMounts)
	containerConfig := createContainerConfig(config.ContainerName,
		config.Image,
		config.Env,
		config.Labels,
		config.Command,
		config.Entrypoint,
		config.VolumeSet)

	return c.ContainerCreate(
		ctx,
		&containerConfig,
		&hostConfig,
		&networkingConfig,
		config.ContainerName,
	)
}

// StartContainer starts a docker container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	options := types.ContainerStartOptions{}
	return c.ContainerStart(ctx, containerID, options)
}

// StopContainer : stop docker container
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	timeout := 60 * time.Second
	return c.ContainerStop(ctx, containerID, &timeout)
}

// RemoveContainer removes a docker container
func (c *Client) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	options := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         force,
	}
	return c.ContainerRemove(ctx, containerID, options)
}

// GetOneContainer returns a docker container if it exists
func (c *Client) GetOneContainer(ctx context.Context, containerId string) (types.ContainerJSON, error) {
	return c.ContainerInspect(ctx, containerId)
}

// GetKraneContainers : gets all containers on the host machine
func (c *Client) GetAllContainers(ctx *context.Context) ([]types.ContainerJSON, error) {
	options := types.ContainerListOptions{
		All:   true,
		Quiet: false,
	}

	containers, err := c.ContainerList(*ctx, options)
	if err != nil {
		return make([]types.ContainerJSON, 0), err
	}

	toJsonContainers := make([]types.ContainerJSON, 0)
	for _, cc := range containers {
		containerJson, err := c.GetOneContainer(*ctx, cc.ID)
		if err != nil {
			return make([]types.ContainerJSON, 0), err
		}

		toJsonContainers = append(toJsonContainers, containerJson)
	}
	return toJsonContainers, nil
}

// GetContainerStatus returns the status of a docker container if it exists
func (c *Client) GetContainerStatus(ctx context.Context, containerID string, stream bool) (stats types.ContainerStats, err error) {
	return c.ContainerStats(ctx, containerID, stream)
}

// StreamContainerLogs streams container logs into ioReader
func (c *Client) StreamContainerLogs(containerID string) (reader io.Reader, err error) {
	ctx := context.Background()
	defer ctx.Done()

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "50",
	}

	return c.ContainerLogs(ctx, containerID, options)
}

// ConnectContainerToNetwork connects a container to a docker network
func (c *Client) ConnectContainerToNetwork(ctx *context.Context, networkID string, containerID string) (err error) {
	config := network.EndpointSettings{NetworkID: networkID}
	return c.NetworkConnect(*ctx, networkID, containerID, &config)
}

func createContainerConfig(
	hostname string,
	image string,
	env []string,
	labels map[string]string,
	command []string,
	entrypoint []string,
	volumes map[string]struct{}) container.Config {
	config := container.Config{
		Hostname: hostname,
		Image:    image,
		Env:      env,
		Labels:   labels,
		Volumes:  volumes,
	}

	if len(command) > 0 {
		config.Cmd = command
	}

	if len(entrypoint) > 0 {
		config.Entrypoint = entrypoint
	}

	return config
}

func createHostConfig(ports nat.PortMap, volumes []mount.Mount) container.HostConfig {
	return container.HostConfig{
		PortBindings: ports,
		AutoRemove:   true,
		Mounts:       volumes,
	}
}
