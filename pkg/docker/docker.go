package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

type DockerClient struct {
	*client.Client
}

var instance *DockerClient
var once sync.Once

// KraneNetworkName : every deployed container will be attached to this network
// TODO: this should be configured somewhere else and passed down when creating docker client / network
var KraneNetworkName = "krane"

// GetClient : get docker client
func GetClient() *DockerClient { return instance }

// Init : starts docker client
func Init() {
	logrus.Info("Connecting to Docker client...")

	once.Do(func() {
		client, err := client.NewEnvClient()
		if err != nil {
			logrus.Fatalf("Failed connecting to Docker client on host machine %s", err.Error())
		}

		instance = &DockerClient{client}

		ctx := context.Background()

		// Create krane docker network
		logrus.Info("Creating Krane Docker network...")
		_, err = instance.CreateBridgeNetwork(&ctx, KraneNetworkName)
		if err != nil {
			logrus.Fatalf("Failed to create Krane Docker network %s", err.Error())
		}

		ctx.Done()
	})
}

// PullImage : poll docker image from registry
func (c *DockerClient) PullImage(ctx *context.Context, image string) (err error) {
	options := types.ImagePullOptions{
		RegistryAuth: "", // RegistryAuth is the base64 encoded credentials for the registry
	}

	reader, err := c.ImagePull(*ctx, image, options)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, reader)
	err = reader.Close()

	return
}

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

// FormatImageSourceURL : format into appropriate docker image url
func FormatImageSourceURL(
	repo string,
	imageName string,
	tag string) string {
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s/%s:%s", repo, imageName, tag)
}

// CreateBridgeNetwork : creates docker bridge network with a given name
func (c *DockerClient) CreateBridgeNetwork(ctx *context.Context, name string) (types.NetworkCreateResponse, error) {
	// Check if krane network already exists
	kNet, err := c.GetNetworkByName(ctx, name)
	if err != nil {
		return types.NetworkCreateResponse{}, err
	}
	if kNet.ID != "" {
		return types.NetworkCreateResponse{ID: kNet.ID}, nil
	}

	// If no exisitng network, create it
	options := types.NetworkCreate{
		Driver:         "bridge",
		CheckDuplicate: true,
	}
	return c.NetworkCreate(*ctx, name, options)
}

// GetNetworkByName : find a netwokr by name on this docker host
func (c *DockerClient) GetNetworkByName(ctx *context.Context, name string) (types.NetworkResource, error) {
	// Get all the networks
	options := types.NetworkListOptions{}
	nets, err := c.NetworkList(*ctx, options)
	if err != nil {
		return types.NetworkResource{}, err
	}

	var kNet types.NetworkResource
	for _, net := range nets {
		if net.Name == name {
			kNet = net
			break
		}
	}

	return kNet, nil
}

// RemoveImage : deletes docker image
func (c *DockerClient) RemoveImage(ctx *context.Context, imageID string) ([]types.ImageDelete, error) {
	options := types.ImageRemoveOptions{
		Force:         false,
		PruneChildren: false,
	}
	return c.ImageRemove(*ctx, imageID, options)
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
