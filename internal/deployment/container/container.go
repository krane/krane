package container

import (
	"context"
	"strings"

	"github.com/biensupernice/krane/docker"
	"github.com/docker/docker/api/types"
)

// Get : container for a deployment
func Get(ctx *context.Context, deploymentName string) (containers []types.Container) {
	containers, err := docker.ListContainers(ctx)
	if err != nil {
		return
	}

	// Remove container not part of this deployment
	// Use the label from the container to determine if its part of the requested deployment
	for i, container := range containers {
		if strings.Compare(deploymentName, container.Labels["deployment.name"]) != 0 {
			containers = append(containers[:i], containers[i+1:]...)
		}
	}

	return
}
