package deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/biensupernice/krane/internal/docker"
)

// KraneContainer represents a Krane managed container
type KraneContainer struct {
	ID         string            `json:"id"`
	Deployment string            `json:"deployment"`
	Name       string            `json:"name"`
	NetworkID  string            `json:"network_id"`
	Image      string            `json:"image"`
	ImageID    string            `json:"image_id"`
	CreatedAt  int64             `json:"created_at"`
	Labels     map[string]string `json:"labels"`
	State      ContainerState    `json:"state"`
	Ports      []Port            `json:"ports"`
	Volumes    []Volume          `json:"volumes"`
	Command    []string          `json:"command"`
	Entrypoint []string          `json:"entrypoint"`
}

// ContainerState represents the state of a Krane container
type ContainerState struct {
	Status     string        `json:"status"` // created,started,running ...
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

type ContainerStatus string

const (
	ContainerRunning ContainerStatus = "running"
	ContainerStarted ContainerStatus = "started"
	ContainerCreated ContainerStatus = "created"
)

// ContainerCreate creates a docker container from a deployment config
func ContainerCreate(config Config) (KraneContainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()

	mappedConfig := config.DockerConfig()
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

// Start starts a Krane managed Docker Container
func (c KraneContainer) Start() error {
	ctx := context.Background()
	defer ctx.Done()

	client := docker.GetClient()
	return client.StartContainer(ctx, c.ID)
}

// Remove removes a Krane managed Docker container
func (c KraneContainer) Remove() error {
	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().RemoveContainer(ctx, c.ID, true)
}

// fromDockerContainerToKcontainer converts a docker container into a KraneContainer
func fromDockerContainerToKcontainer(container types.ContainerJSON) KraneContainer {
	ctx := context.Background()
	defer ctx.Done()

	createdAt, _ := time.Parse(time.RFC3339, container.ContainerJSONBase.Created)
	state := fromDockerStateToState(*container.State)
	ports := fromPortMapToPortList(container.NetworkSettings.Ports)
	volumes := fromMountPointToVolumeList(container.Mounts)

	return KraneContainer{
		ID:         container.ID,
		Deployment: container.Config.Labels[docker.ContainerNamespaceLabel],
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

func fromDockerStateToState(state types.ContainerState) ContainerState {
	return ContainerState{
		Status:     state.Status,
		Running:    state.Running,
		Paused:     state.Paused,
		Restarting: state.Restarting,
		OOMKilled:  state.OOMKilled,
		Dead:       state.Dead,
		Pid:        state.Pid,
		ExitCode:   state.ExitCode,
		Error:      state.Error,
		StartedAt:  state.StartedAt,
		FinishedAt: state.FinishedAt,
		Health:     state.Health,
	}
}

// GetContainers get all containers managed by Krane
func GetContainers() ([]KraneContainer, error) {
	ctx := context.Background()
	defer ctx.Done()

	allContainers, err := docker.GetClient().GetAllContainers(&ctx)
	if err != nil {
		return make([]KraneContainer, 0), err
	}

	// filter for Krane managed containers
	containers := make([]KraneContainer, 0)
	for _, container := range allContainers {
		if isKraneManagedContainer(container) {
			containers = append(containers, fromDockerContainerToKcontainer(container))
		}
	}

	return containers, nil
}

// isKraneManagedContainer returns if a container is managed by Krane based on its labels
func isKraneManagedContainer(container types.ContainerJSON) bool {
	return len(container.Config.Labels[docker.ContainerNamespaceLabel]) > 0
}

// GetContainersByDeployment get containers filtered by deployment
func GetContainersByDeployment(deployment string) ([]KraneContainer, error) {
	allContainers, err := GetContainers()
	if err != nil {
		return make([]KraneContainer, 0), err
	}

	// filter by deployment
	containers := make([]KraneContainer, 0)
	for _, container := range allContainers {
		if deployment == container.Deployment {
			containers = append(containers, container)
		}
	}

	return containers, nil
}

// RetriableContainerHealthCheck returns an error if a container is considered unhealthy
func RetriableContainerHealthCheck(containers []KraneContainer, pollRetry int) error {
	for _, c := range containers {
		for i := 0; i <= pollRetry; i++ {
			expBackOff := time.Duration(10 * i)
			time.Sleep(expBackOff * time.Second)

			isRunning, err := c.Running()
			if err != nil {
				if i == pollRetry {
					return fmt.Errorf("container is not healthy %v", err)
				}
				continue
			}

			if !isRunning {
				if i == pollRetry {
					return fmt.Errorf("container is not healthy %v", err)
				}
				continue
			}

			// if reached here container in a running state
			break
		}
	}
	return nil
}

// Running returns whether a container is in a running state
func (c KraneContainer) Running() (bool, error) {
	ctx := context.Background()
	defer ctx.Done()

	resp, err := docker.GetClient().GetOneContainer(ctx, c.ID)
	if err != nil {
		return false, err
	}

	if resp.State.Running {
		return true, nil
	}

	return false, fmt.Errorf("container %s is not in running state", c.ID)

}
