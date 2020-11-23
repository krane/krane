package container

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/lithammer/shortuuid/v3"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
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
	Command   []string          `json:"command"`
}

type State struct {
	Status     string        `json:"status"`
	Running    bool          `json:"running"`
	Paused     bool          `json:"paused"`
	Restarting bool          `json:"restarting"`
	OOMKilled  bool          `json:"oom_killed"`
	Dead       bool          `json:"dead"`
	Pid        int           `json:"pid"`
	ExitCode   int           `json:"exit_code"`
	Error      string        `json:"error"`
	StartedAt  string        `json:"started"`
	FinishedAt string        `json:"finished_at"`
	Health     *types.Health `json:",omitempty"`
}

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

	return mapContainerJsonToKcontainer(json), nil
}

func mapContainerJsonToKcontainer(container types.ContainerJSON) Kcontainer {
	createdAt, _ := time.Parse(time.RFC3339, container.Created)
	envs := fromDockerToEnvMap(container.Config.Env)

	return Kcontainer{
		ID:        container.ID,
		Name:      container.Name,
		NetworkID: container.NetworkSettings.EndpointID,
		Image:     container.Config.Image,
		ImageID:   container.Image,
		Env:       envs,
		CreatedAt: createdAt.Unix(),
		Command:   container.Config.Cmd,
		Labels:    container.Config.Labels,
		// TODO: resolve rest of the fields
	}
}

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

	containerName := fmt.Sprintf("%s-%s", cfg.Name, shortuuid.New())
	return docker.CreateContainerConfig{
		ContainerName: containerName,
		Image:         cfg.Image,
		NetworkID:     knetwork.ID,
		Labels:        labels,
		Ports:         ports,
		Volumes:       volumes,
		Env:           envars,
		Command:       cfg.Command,
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
func fromDockerContainerToKcontainer(container types.ContainerJSON) Kcontainer {
	ctx := context.Background()
	defer ctx.Done()

	createdAt, _ := strconv.ParseInt(container.ContainerJSONBase.Created, 10, 64)
	ports := fromDockerToKconfigPortMap(container.NetworkSettings.Ports)

	// volumes

	// Map Env

	// Map Secrets

	return Kcontainer{
		ID:        container.ID,
		Namespace: container.Config.Labels[KraneContainerNamespaceLabel],
		Name:      container.Name,
		Image:     container.Config.Image,
		ImageID:   container.ContainerJSONBase.Image,
		CreatedAt: createdAt,
		Labels:    container.Config.Labels,
		NetworkID: container.NetworkSettings.Networks[docker.KraneNetworkName].NetworkID,
		// State:     container.ContainerJSONBase.State,
		Status:  container.ContainerJSONBase.State.Status,
		Ports:   ports,
		Command: container.Config.Cmd,
		// Env:
		// Volumes:
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

	kcontainers := make([]Kcontainer, 0)
	for _, container := range containers {
		if isKraneManagedContainer(container) {
			kcontainers = append(kcontainers, fromDockerContainerToKcontainer(container))
		}
	}

	return kcontainers, nil
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

func isKraneManagedContainer(container types.ContainerJSON) bool {
	namespaceLabel := container.Config.Labels[KraneContainerNamespaceLabel]
	if namespaceLabel == "" {
		return false
	}
	return true
}
