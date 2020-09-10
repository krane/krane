package scheduler

import (
	"github.com/docker/docker/api/types"

	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/kranecfg"
)

type stateMapping struct {
	desiredState kranecfg.KraneConfig
	containers   []types.Container
}

func mapDeploymentsToContainers(deployments []kranecfg.KraneConfig, containers []types.Container) map[string]stateMapping {
	m := make(map[string]stateMapping, 0)

	for _, d := range deployments {
		deploymentContainers := docker.FilterContainersByDeployment(d.Name, containers)
		m[d.Name] = stateMapping{d, deploymentContainers}
	}

	return m
}
