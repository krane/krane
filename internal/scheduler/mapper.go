package scheduler

import (
	"encoding/json"

	"github.com/docker/docker/api/types"

	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/kranecfg"
)

type DeploymentToContainersMapping struct {
	config     kranecfg.KraneConfig
	containers []types.Container
}

func mapDeploymentsToContainers(deployments []kranecfg.KraneConfig, containers []types.Container) map[string]DeploymentToContainersMapping {
	m := make(map[string]DeploymentToContainersMapping, 0)

	for _, d := range deployments {
		deploymentContainers := docker.FilterContainersByDeployment(d.Name, containers)
		m[d.Name] = DeploymentToContainersMapping{d, deploymentContainers}
	}

	return m
}

func (s *DeploymentToContainersMapping) serialize() ([]byte, error) { return json.Marshal(s) }
