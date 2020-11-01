package container

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/lithammer/shortuuid/v3"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/secrets"
)

// Krane custom container struct
type Kcontainer struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	NetworkID string            `json:"network_id"`
	Image     string            `json:"image"`
	ImageID   string            `json:"image_id"`
	CreatedAt int64             `json:"created_at"`
	Labels    map[string]string `json:"labels"`
	State     State             `json:"state"`  // ex: running
	Status    string            `json:"status"` // ex: Up 17 hours
	Ports     []Port            `json:"ports"`
	Volumes   []Volume          `json:"volumes"`
	Env       map[string]string `json:"env"`
	Secrets   []secrets.Secret  `json:"secrets"`
	Command   string            `json:"command"`
}

type State string

const (
	Running State = "running"
	Unknown State = "unknown"
)

func Create(cfg config.Kconfig) (Kcontainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()

	mappedConfig := fromKconfigToCreateContainerConfig(cfg)
	body, err := client.CreateContainer(ctx, mappedConfig)
	if err != nil {
		return Kcontainer{}, err
	}

	// get container
	json, err := client.GetOneContainer(ctx, body.ID)
	if err != nil {
		return Kcontainer{}, err
	}

	return mapContainerJsonToKontainer(json), nil
}

func mapContainerJsonToKontainer(container types.ContainerJSON) Kcontainer {
	envs := fromDockerToEnvMap(container.Config.Env)
	return Kcontainer{
		ID:        container.ID,
		Name:      container.Name,
		Namespace: container.Name,
		NetworkID: container.NetworkSettings.EndpointID,
		Image:     container.Image,
		ImageID:   container.Image,
		Env:       envs,
		// TODO: resolve rest of the fields
	}
}

func fromKconfigToCreateContainerConfig(cfg config.Kconfig) docker.CreateContainerConfig {
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

	containerName := fmt.Sprintf("%s-%s", cfg.Name, shortuuid.New())
	return docker.CreateContainerConfig{
		ContainerName: containerName,
		Image:         cfg.Image,
		NetworkID:     knetwork.ID,
		Labels:        labels,
		Ports:         ports,
		Volumes:       volumes,
		Env:           envars,
		Cmd:           cfg.Command,
	}
}

func (k Kcontainer) Start() error {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	return client.StartContainer(ctx, k.ID)
}

func (k Kcontainer) Stop() error {
	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().StopContainer(ctx, k.ID)
}

func (k Kcontainer) Remove() error {
	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().RemoveContainer(ctx, k.ID, true)
}

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

// convert docker container into a Kcontainer
func fromDockerContainerToKcontainer(container types.Container) Kcontainer {
	// Map container state
	var containerState State
	switch container.State {
	case "running":
		containerState = Running
	default:
		containerState = Unknown
	}

	// Map ports
	ports := fromDockerToKcontainerPorts(container.Ports)

	// Map Env

	// Map Secrets

	return Kcontainer{
		ID:        container.ID,
		Namespace: container.Labels[KraneContainerNamespaceLabel],
		Name:      container.Names[0],
		Image:     container.Image,
		ImageID:   container.ImageID,
		CreatedAt: container.Created,
		Labels:    container.Labels,
		NetworkID: container.NetworkSettings.Networks[docker.KraneNetworkName].NetworkID,
		State:     containerState,
		Status:    container.Status,
		Ports:     ports,
		Command:   container.Command,
		// TODO: the rest
	}
}

// get krane managed docker containers mapped to a Kcontainer
func GetAllContainers(client *docker.Client) ([]Kcontainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	containers, err := client.GetAllContainers(&ctx)
	if err != nil {
		return make([]Kcontainer, 0), err
	}

	kontainers := make([]Kcontainer, 0)
	for _, container := range containers {
		// For krane managed containers map to a Kcontainer
		if isKraneManagedContainer(container) {
			kontainers = append(kontainers, fromDockerContainerToKcontainer(container))
		}
	}

	return kontainers, nil
}

// filter krane manage containers by namespace
func GetKontainersByNamespace(namespace string) ([]Kcontainer, error) {
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

// Container label for krane managed containers
const KraneContainerNamespaceLabel = "krane.deployment.namespace"

func isKraneManagedContainer(container types.Container) bool {
	namespaceLabel := container.Labels[KraneContainerNamespaceLabel]
	if namespaceLabel == "" {
		return false
	}

	return true
}
