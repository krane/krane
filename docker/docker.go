package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	dkrClient *client.Client // Single docker client
)

// New : create docker client
func New() (*client.Client, error) {
	if dkrClient != nil {
		return dkrClient, nil
	}

	client, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	dkrClient = client

	return dkrClient, nil
}

// PullImage : poll docker image from registry
func PullImage(ctx *context.Context, image string) error {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return err
	}

	options := types.ImagePullOptions{
		RegistryAuth: "", // RegistryAuth is the base64 encoded credentials for the registry
	}
	ioreader, err := dkrClient.ImagePull(*ctx, image, options)

	if err != nil {
		return err
	}

	io.Copy(os.Stdout, ioreader)
	err = ioreader.Close()
	if err != nil {
		return err
	}

	return nil
}

// CreateContainer blah
func CreateContainer(
	ctx *context.Context,
	image string,
	containerName string,
	hPort string,
	cPort string,
) (container.ContainerCreateCreatedBody, error) {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return container.ContainerCreateCreatedBody{}, err
	}

	// Configure Host Port
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: hPort,
	}

	// Configure Container Port
	containerPort, err := nat.NewPort("tcp", cPort)
	if err != nil {
		log.Printf("Unable to configure container port %s - %s", cPort, err.Error())
		return container.ContainerCreateCreatedBody{}, err
	}

	// Bind host-to-container ports
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	// Setup host conf
	hostConf := &container.HostConfig{PortBindings: portBinding}

	// Setup container conf
	containerConf := &container.Config{
		Hostname: "blah",
		Image:    image,
		Env:      []string{"TEST_ENV=pipi"},
		Labels:   map[string]string{"TEST_LABEL": "poopoo"},
	}

	// Setup networking conf
	networkConf := &network.NetworkingConfig{}

	return dkrClient.ContainerCreate(*ctx, containerConf, hostConf, networkConf, containerName)
}

// StartContainer blah
func StartContainer(ctx *context.Context, containerID string) error {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return err
	}

	options := types.ContainerStartOptions{}
	return dkrClient.ContainerStart(*ctx, containerID, options)
}

// StopContainer : stop docker container
func StopContainer(ctx *context.Context, containerID string) error {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return err
	}

	return dkrClient.ContainerStop(*ctx, containerID, nil)
}

// RemoveContainer : remove docker container
func RemoveContainer(ctx *context.Context, containerID string) error {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return err
	}

	options := types.ContainerRemoveOptions{}
	return dkrClient.ContainerRemove(*ctx, containerID, options)
}

// ListContainers : get all containers
func ListContainers(ctx *context.Context) (containers []types.Container, err error) {
	if dkrClient == nil {
		err = fmt.Errorf("docker client not initialized")
		return
	}
	options := types.ContainerListOptions{}
	return dkrClient.ContainerList(*ctx, options)
}

// FormatImageSourceURL : format into appropriate docker image url
func FormatImageSourceURL(
	repo string,
	imageName string,
	tag string) string {
	return fmt.Sprintf("%s/%s:%s", repo, imageName, tag)
}

// Helper to find the current host ip address - 0.0.0.0 binds to all ip's
func getHostIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return "0.0.0.0"
}
