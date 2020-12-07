package container

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/lithammer/shortuuid/v3"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/docker"
)

// KraneContainer : custom container representation for Krane managed containers
type KraneContainer struct {
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

// Create : create docker container
func Create(cfg config.DeploymentConfig) (KraneContainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()

	// TODO: add these methods to the deployment config struct
	mappedConfig := fromKconfigToCreateContainerConfig(cfg)
	body, err := docker.GetClient().CreateContainer(ctx, mappedConfig)
	if err != nil {
		return KraneContainer{}, err
	}

	// the response from creating a container doesnt provide enough information
	// about the resources it created, we need to inspect the containers for full details
	json, err := client.GetOneContainer(ctx, body.ID)
	if err != nil {
		return KraneContainer{}, err
	}

	return fromDockerContainerToKcontainer(json), nil
}

// Start : start a KraneContainer
func (c KraneContainer) Start() error {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	return client.StartContainer(ctx, c.ID)
}

// Stop : stop a KraneContainer
func (c KraneContainer) Stop() error {
	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().StopContainer(ctx, c.ID)
}

// Remove : remove a KraneContainer
func (c KraneContainer) Remove() error {
	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().RemoveContainer(ctx, c.ID, true)
}

// Ok : checks if the container is in a running state
func (c KraneContainer) Ok() (bool, error) {
	ctx := context.Background()
	defer ctx.Done()

	resp, err := docker.GetClient().GetOneContainer(ctx, c.ID)
	if err != nil {
		return false, err
	}

	if !resp.State.Running {
		return false, fmt.Errorf("container %s is not in running state", c.ID)
	}

	return true, nil
}

func (c KraneContainer) toContainer() types.Container { return types.Container{} }

// GetAllContainers : get all containers
func GetAllContainers(client *docker.Client) ([]KraneContainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	containers, err := client.GetAllContainers(&ctx)
	if err != nil {
		return make([]KraneContainer, 0), err
	}

	// filter for Krane managed containers
	kcontainers := make([]KraneContainer, 0)
	for _, container := range containers {
		if isKraneManagedContainer(container) {
			kcontainers = append(kcontainers, fromDockerContainerToKcontainer(container))
		}
	}

	return kcontainers, nil
}

// GetContainersByDeployment : get containers filtered by namespace
func GetContainersByDeployment(namespace string) ([]KraneContainer, error) {
	allContainers, err := GetAllContainers(docker.GetClient())
	if err != nil {
		return make([]KraneContainer, 0), err
	}

	// filter by deployment
	containers := make([]KraneContainer, 0)
	for _, container := range allContainers {
		if namespace == container.Namespace {
			containers = append(containers, container)
		}
	}

	return containers, nil
}

// isKraneManagedContainer : check if a container is a Krane managed container
func isKraneManagedContainer(container types.ContainerJSON) bool {
	return len(container.Config.Labels[KraneContainerLabel]) > 0
}

// fromKconfigToCreateContainerConfig :
func fromKconfigToCreateContainerConfig(cfg config.DeploymentConfig) docker.CreateContainerConfig {
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

// fromDockerContainerToKcontainer : convert docker container into a KraneContainer
func fromDockerContainerToKcontainer(container types.ContainerJSON) KraneContainer {
	ctx := context.Background()
	defer ctx.Done()

	createdAt, _ := time.Parse(time.RFC3339, container.ContainerJSONBase.Created)
	state := fromDockerStateToKstate(*container.State)
	ports := fromDockerToKconfigPortMap(container.NetworkSettings.Ports)
	volumes := fromMountPointToKconfigVolumes(container.Mounts)

	return KraneContainer{
		ID:         container.ID,
		Namespace:  container.Config.Labels[KraneContainerLabel],
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
