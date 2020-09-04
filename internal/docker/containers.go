package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

var (
	// KraneContainerLabelName : used to identify krane managed containers
	KraneContainerLabelName = "krane.deployment.name"
)

func GetContainers(ctx *context.Context, deploymentName string) ([]types.Container, error) {

	client := GetClient()

	// Find all containers
	allContainers, err := client.GetAllContainers(ctx)
	if err != nil {
		return make([]types.Container, 0), err
	}

	deploymentContainers := make([]types.Container, 0)
	for _, container := range allContainers {
		kraneLabel := container.Labels[KraneContainerLabelName]
		if kraneLabel == deploymentName {
			deploymentContainers = append(deploymentContainers, container)
		}
	}

	return deploymentContainers, nil
}