package container

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/lithammer/shortuuid/v3"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/secrets"
	"github.com/biensupernice/krane/internal/proxy"
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
	body, err := client.CreateContainer(ctx, mappedConfig)
	if err != nil {
		return Kontainer{}, err
	}

	// get container
	json, err := client.GetOneContainer(ctx, body.ID)
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

	knetwork, err := docker.GetClient().GetNetworkByName(ctx, docker.KraneNetworkName)
	if err != nil {
		return docker.CreateContainerConfig{}
	}

	envars := makeContainerEnvars(cfg)
	labels := makeContainerLabels(cfg)
	volumes := makeContainerVolumes(cfg)
	ports := makeContainerPorts(cfg)

	containerName := fmt.Sprintf("%s-%s", cfg.Name, shortuuid.New())
	return docker.CreateContainerConfig{
		ContainerName: containerName,
		Image:         cfg.Image,
		NetworkID:     knetwork.ID,
		Labels:        labels,
		Ports:         ports,
		Volumes:       volumes,
		Env:           envars,
	}
}

func makeContainerPorts(cfg config.Config) nat.PortMap {
	bindings := nat.PortMap{}
	for hostPort, containerPort := range cfg.Ports {
		// host port
		hostBinding := nat.PortBinding{HostPort: hostPort}

		// container port
		// TODO: figure out if we can bind ports of other types besides tcp
		cPort, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			continue
		}

		bindings[cPort] = []nat.PortBinding{hostBinding}
	}

	return bindings
}

func makeContainerVolumes(cfg config.Config) []mount.Mount {
	volumes := make([]mount.Mount, 0)

	for hostVolume, containerVolume := range cfg.Volumes {
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: hostVolume,
			Target: containerVolume,
		})
	}

	return volumes
}

func makeContainerLabels(cfg config.Config) map[string]string {
	labels := map[string]string{
		KraneContainerNamespaceLabel: cfg.Name,
	}

	// TODO: theres a bug where it only applies a single label if aliases is greater than 1.
	// This is because the labels key get overwritten
	for _, alias := range cfg.Alias {
		routingLabels := proxy.MakeContainerRoutingLabels(cfg.Name, alias)
		for _, label := range routingLabels {
			labels[label.Label] = label.Value
		}
	}

	return labels
}

// convert map of envars into formatted list of envars
func makeContainerEnvars(cfg config.Config) []string {
	envs := make([]string, 0)

	// config environment variables
	for k, v := range cfg.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	// resolve secret by alias
	for key, alias := range cfg.Secrets {
		secret, err := secrets.Get(cfg.Name, alias)
		if err != nil || secret == nil {
			logrus.Debugf("Unable to get resolve secret for %s with alias @%s", cfg.Name, alias)
			continue
		}
		envs = append(envs, fmt.Sprintf("%s=%s", key, secret.Value))
	}

	return envs
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

	return docker.GetClient().RemoveContainer(ctx, k.ID, true)
}

func (k Kontainer) Ok() (bool, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	status, err := client.GetOneContainer(ctx, k.ID)
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
func GetKontainersByNamespace(namespace string) ([]Kontainer, error) {
	client := docker.GetClient()
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
