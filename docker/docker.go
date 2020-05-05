package docker

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// NewClient blah
func NewClient() (*client.Client, error) {
	return client.NewEnvClient()
}

// PullImage blah
func PullImage(
	ctx *context.Context,
	dockerClient *client.Client,
	image string) error {
	ioreader, err := dockerClient.ImagePull(*ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
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
	dockerClient *client.Client,
	image string,
	hPort string,
	cPort string,
) (container.ContainerCreateCreatedBody, error) {
	// Configure Host Port
	hostBinding := nat.PortBinding{
		HostIP:   getHostIP(),
		HostPort: hPort,
	}

	// Configure Container Port
	containerPort, err := nat.NewPort("tcp", cPort)
	if err != nil {
		panic("Unable to get the port")
	}

	// Bind Host--Container posrts
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	containerConf := &container.Config{Image: image}
	hostConf := &container.HostConfig{PortBindings: portBinding}
	return dockerClient.ContainerCreate(*ctx, containerConf, hostConf, nil, "")
}

// StartContainer blah
func StartContainer(
	ctx *context.Context,
	dockerClient *client.Client,
	containerID string) error {
	return dockerClient.ContainerStart(*ctx, containerID, types.ContainerStartOptions{})
}

// StopContainer blah
func StopContainer(
	ctx *context.Context,
	dockerClient *client.Client,
	containerID string,
) error {
	return dockerClient.ContainerStop(*ctx, containerID, nil)
}

// GetDockerImageSource blah
func FormatImageSourceUrl(
	repo string,
	imageName string,
	tag string) string {
	return fmt.Sprintf("%s/%s:%s", repo, imageName, tag)
}

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
