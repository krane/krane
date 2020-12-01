package container

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/lithammer/shortuuid/v3"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/docker"
)

// Kcontainer : custom container struct for Krane managed containers
type Kcontainer struct {
	ID         string            `json:"id"`
	Namespace  string            `json:"namespace"`
	Name       string            `json:"name"`
	NetworkID  string            `json:"network_id"`
	Image      string            `json:"image"`
	ImageID    string            `json:"image_id"`
	CreatedAt  int64             `json:"created_at"`
	Labels     map[string]string `json:"labels"`
	State      State             `json:"state"`
	Ports      []Port            `json:"ports"`
	Volumes    []Volume          `json:"volumes"`
	Command    []string          `json:"command"`
	Entrypoint []string          `json:"entrypoint"`
}

// Create : create docker container from Kconfig
func Create(cfg kconfig.Kconfig) (Kcontainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()

	mappedConfig := fromKconfigToCreateContainerConfig(cfg)
	body, err := docker.GetClient().CreateContainer(ctx, mappedConfig)
	if err != nil {
		return Kcontainer{}, err
	}

	// the response from creating a container doesnt provide enough information
	// about the resources it created, we need to inspect the containers for full details
	json, err := client.GetOneContainer(ctx, body.ID)
	if err != nil {
		return Kcontainer{}, err
	}

	return fromDockerContainerToKcontainer(json), nil
}

// Start : start Kcontainer
func (k Kcontainer) Start() error {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	return client.StartContainer(ctx, k.ID)
}

// Stop : stop Kcontainer
func (k Kcontainer) Stop() error {
	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().StopContainer(ctx, k.ID)
}

// Remove : remove Kcontainer
func (k Kcontainer) Remove() error {
	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().RemoveContainer(ctx, k.ID, true)
}

// Ok : returns if container is in a running state
func (k Kcontainer) Ok() (bool, error) {
	ctx := context.Background()
	defer ctx.Done()

	status, err := docker.GetClient().GetOneContainer(ctx, k.ID)
	if err != nil {
		return false, err
	}

	if !status.State.Running {
		return false, errors.New("container not in running state")
	}

	return true, nil
}

func (k Kcontainer) toContainer() types.Container { return types.Container{} }

// GetAllContainers : get all containers as Kcontainers
func GetAllContainers(client *docker.Client) ([]Kcontainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	containers, err := client.GetAllContainers(&ctx)
	if err != nil {
		return make([]Kcontainer, 0), err
	}

	// filter krane managed containers
	kcontainers := make([]Kcontainer, 0)
	for _, container := range containers {
		if isKraneManagedContainer(container) {
			kcontainers = append(kcontainers, fromDockerContainerToKcontainer(container))
		}
	}

	return kcontainers, nil
}

// GetContainersByNamespace : get Kcontainers filtered by namespace
func GetContainersByNamespace(namespace string) ([]Kcontainer, error) {
	client := docker.GetClient()

	// get all containers managed by krane
	containers, err := GetAllContainers(client)
	if err != nil {
		return make([]Kcontainer, 0), err
	}

	// filter containers for just this deployment
	filteredKontainers := make([]Kcontainer, 0)
	for _, containers := range containers {
		if namespace == containers.Namespace {
			filteredKontainers = append(filteredKontainers, containers)
		}
	}

	return filteredKontainers, nil
}

// isKraneManagedContainer : check if a container is a Krane managed container
func isKraneManagedContainer(container types.ContainerJSON) bool {
	namespaceLabel := container.Config.Labels[KraneContainerNamespaceLabel]
	if namespaceLabel == "" {
		return false
	}
	return true
}

// fromKconfigToCreateContainerConfig :
func fromKconfigToCreateContainerConfig(cfg kconfig.Kconfig) docker.CreateContainerConfig {
	ctx := context.Background()
	defer ctx.Done()

	knetwork, err := docker.GetClient().GetNetworkByName(ctx, docker.KraneNetworkName)
	if err != nil {
		return docker.CreateContainerConfig{}
	}

	envars := fromKconfigDockerEnvList(cfg)
	labels := fromKconfigToDockerLabelMap(cfg)
	volumes := fromKconfigToDockerVolumeMount(cfg)
	ports := fromKconfigToDockerPortMap(cfg)

	var command []string
	var entrypoint []string

	if cfg.Command != "" {
		command = append(command, cfg.Command)
	}

	if cfg.Entrypoint != "" {
		entrypoint = append(entrypoint, cfg.Entrypoint)
	}

	containerName := fmt.Sprintf("%s-%s", cfg.Name, shortuuid.New())
	return docker.CreateContainerConfig{
		ContainerName: containerName,
		Image:         cfg.Image,
		NetworkID:     knetwork.ID,
		Labels:        labels,
		Ports:         ports,
		Volumes:       volumes,
		Env:           envars,
		Command:       command,
		Entrypoint:    entrypoint,
	}
}

// fromDockerContainerToKcontainer : convert docker container into a Kcontainer
func fromDockerContainerToKcontainer(container types.ContainerJSON) Kcontainer {
	ctx := context.Background()
	defer ctx.Done()

	createdAt, _ := time.Parse(time.RFC3339, container.ContainerJSONBase.Created)
	state := fromDockerStateToKstate(*container.State)
	ports := fromDockerToKconfigPortMap(container.NetworkSettings.Ports)
	volumes := fromMountPointToKconfigVolumes(container.Mounts)

	return Kcontainer{
		ID:         container.ID,
		Namespace:  container.Config.Labels[KraneContainerNamespaceLabel],
		Name:       container.Name,
		NetworkID:  container.NetworkSettings.Networks[docker.KraneNetworkName].NetworkID,
		Image:      container.Config.Image,
		ImageID:    container.ContainerJSONBase.Image,
		CreatedAt:  createdAt.Unix(),
		Labels:     container.Config.Labels,
		State:      state,
		Ports:      ports,
		Volumes:    volumes,
		Command:    container.Config.Cmd,
		Entrypoint: container.Config.Entrypoint,
	}
}
