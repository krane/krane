package docker

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"

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

// CreateContainer : create docker container
func CreateContainer(
	ctx *context.Context,
	image string,
	containerName string,
	networkID string,
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
		return container.ContainerCreateCreatedBody{}, err
	}

	// Bind host-to-container ports
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	// Setup host conf
	hostConf := &container.HostConfig{
		PortBindings: portBinding,
		AutoRemove:   false,
	}

	// Setup container conf
	containerConf := &container.Config{
		Hostname: "bsn",
		Image:    image,
		Env:      []string{"TEST_ENV=pipi"},
		Labels:   map[string]string{"TEST_LABEL": "poopoo"},
	}

	// Setup networking conf
	endpointConf := map[string]*network.EndpointSettings{"krane": &network.EndpointSettings{NetworkID: networkID}}
	networkConf := &network.NetworkingConfig{EndpointsConfig: endpointConf}

	return dkrClient.ContainerCreate(*ctx, containerConf, hostConf, networkConf, containerName)
}

// StartContainer : start docker container
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

	timeout := 60 * time.Second
	return dkrClient.ContainerStop(*ctx, containerID, &timeout)
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

// CreateBridgeNetwork : creates docker bridge network with a given name
func CreateBridgeNetwork(ctx *context.Context, name string) (types.NetworkCreateResponse, error) {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return types.NetworkCreateResponse{}, err
	}

	// Check if krane network already exists
	kNet, err := GetNetworkByName(ctx, name)
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
	return dkrClient.NetworkCreate(*ctx, name, options)
}

// GetNetworkByName : find a netwokr by name on this docker host
func GetNetworkByName(ctx *context.Context, name string) (types.NetworkResource, error) {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return types.NetworkResource{}, err
	}

	// Get all the networks
	options := types.NetworkListOptions{}
	nets, err := dkrClient.NetworkList(*ctx, options)
	if err != nil {
		return types.NetworkResource{}, err
	}

	var kNet types.NetworkResource
	for i := 0; i < len(nets); i++ {
		if nets[i].Name == name {
			kNet = nets[i]
			break
		}
	}

	return kNet, nil
}

// RemoveImage : deletes docker image
func RemoveImage(ctx *context.Context, imageID string) ([]types.ImageDelete, error) {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return []types.ImageDelete{}, err
	}

	options := types.ImageRemoveOptions{
		Force:         false,
		PruneChildren: false,
	}
	return dkrClient.ImageRemove(*ctx, imageID, options)
}

// ReadContainerLogs :
func ReadContainerLogs(ctx *context.Context, containerID string) (reader io.Reader, err error) {
	if dkrClient == nil {
		err := fmt.Errorf("docker client not initialized")
		return nil, err
	}

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "50",
	}

	return dkrClient.ContainerLogs(*ctx, containerID, options)
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
