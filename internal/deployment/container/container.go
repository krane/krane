package container

import (
	"context"

	"github.com/docker/docker/api/types"

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
	return client.RemoveContainer(ctx, k.ID)
}

func create(cfg config.Config) (Kontainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	// TODO: figure out what to return
	body, err := client.CreateContainer2(ctx, k)
	if err != nil {
		return Kontainer{}, err
	}

	// get container
	_, err = client.GetOneContainer(&ctx, body.ID)
	if err != nil {
		return Kontainer{}, err
	}

	return Kontainer{}, err
}
func (k Kontainer) toContainer() types.Container { return types.Container{} }

// convert docker container into a Kontainer
func fromContainer(container types.Container) Kontainer { return Kontainer{} }

// get krane managed containers
func GetKontainers(client *docker.DockerClient) ([]Kontainer, error) {
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
func GetKontainersByNamespace(client *docker.DockerClient, namespace string) ([]Kontainer, error) {
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
