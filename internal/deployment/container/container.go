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
type Kontainer struct {
	ID        string            `json:"id"`
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	NetworkID string            `json:"network_id"`
	Image     string            `json:"image"`
	ImageID   string            `json:"image_id"`
	CreatedAt int64             `json:"created_at"`
	Labels    map[string]string `json:"labels"`
	State     string            `json:"state"`
	Status    string            `json:"status"`
	Ports     []Port            `json:"ports"`
	Volumes   []Volume          `json:"volumes"`
	Env       map[string]string `json:"env"`
	Secrets   []secrets.Secret  `json:"secrets"`
}

func fromConfig(config config.Config) Kontainer {
	return Kontainer{}
}

func Create(cfg config.Config) (Kontainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()

	mappedConfig := mapConfigToCreateContainerConfig(cfg)
	body, err := client.CreateContainer2(ctx, mappedConfig)
	if err != nil {
		return Kontainer{}, err
	}

	// get container
	json, err := client.GetOneContainer(&ctx, body.ID)
	if err != nil {
		return Kontainer{}, err
	}

	return mapContainerJsonToKontainer(json), nil
}

func mapContainerJsonToKontainer(container types.ContainerJSON) Kontainer {
	return Kontainer{
		ID:        container.ID,
		Name:      container.Name,
		Namespace: container.Name,
		NetworkID: container.NetworkSettings.EndpointID,
		Image:     container.Image,
		ImageID:   container.Image,
		// TODO: resolve rest of the fields
	}
}

func mapConfigToCreateContainerConfig(cfg config.Config) docker.CreateContainerConfig {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	ntw, err := client.GetNetworkByName(ctx, docker.KraneNetworkName)
	if err != nil {
		return docker.CreateContainerConfig{}
	}

	containerName := fmt.Sprintf("%s-%s", cfg.Name, shortuuid.New())
	return docker.CreateContainerConfig{
		ContainerName: containerName,
		Image:         cfg.Image,
		NetworkID:     ntw.ID,
		Labels:        map[string]string{KraneContainerNamespaceLabel: cfg.Name},
		Ports:         nil,
		Volumes:       nil,
		Env:           []string{},
	}
}

func (k Kontainer) Start() error {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	return client.StartContainer(ctx, k.ID)
}

func (k Kontainer) Stop() error {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	return client.StopContainer(ctx, k.ID)
}

func (k Kontainer) Remove() error {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()

	if err := client.RemoveContainer(ctx, k.ID, true); err != nil {
		return err
	}

	return nil
}

func (k Kontainer) Ok() (bool, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	status, err := client.GetOneContainer(&ctx, k.ID)
	if err != nil {
		return false, err
	}

	if !status.State.Running {
		return false, errors.New("container not in running state")
	}

	return true, nil
}

func (k Kontainer) toContainer() types.Container { return types.Container{} }

// convert docker container into a Kontainer
func fromContainer(container types.Container) Kontainer {
	return Kontainer{
		ID:        container.ID,
		Namespace: container.Labels[KraneContainerNamespaceLabel],
		Name:      container.Names[0],
		Image:     container.Image,
		ImageID:   container.ImageID,
		CreatedAt: container.Created,
		Labels:    container.Labels,
		NetworkID: container.NetworkSettings.Networks[docker.KraneNetworkName].NetworkID,
		// TODO: the rest
	}
}

// get krane managed containers
func GetKontainers(client *docker.Client) ([]Kontainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	containers, err := client.GetAllContainers(&ctx)
	if err != nil {
		return make([]Kontainer, 0), err
	}

	kontainers := make([]Kontainer, 0)
	for _, container := range containers {
		if isKraneManagedContainer(container) {
			k := fromContainer(container)
			kontainers = append(kontainers, k)
		}
	}

	return kontainers, nil
}

// filter krane manage containers by namespace
func GetKontainersByNamespace(client *docker.Client, namespace string) ([]Kontainer, error) {
	kontainers, err := GetKontainers(client)
	if err != nil {
		return make([]Kontainer, 0), err
	}

	filteredKontainers := make([]Kontainer, 0)
	for _, kontainer := range kontainers {
		if namespace == kontainer.Namespace {
			filteredKontainers = append(filteredKontainers, kontainer)
		}
	}

	return filteredKontainers, nil
}

// Container label for identifying krane managed containers
const KraneContainerNamespaceLabel = "krane.deployment.namespace"

func isKraneManagedContainer(container types.Container) bool {
	namespaceLabel := container.Labels[KraneContainerNamespaceLabel]
	if namespaceLabel == "" {
		return false
	}

	return true
}