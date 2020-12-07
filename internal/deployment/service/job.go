package service

import (
	"fmt"

	"github.com/docker/distribution/uuid"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/utils"
)

type action string

// Follows docker conventions
// Up : to create container resources
// Down : to remove container resources
const (
	Up   action = "UP"
	Down action = "DOWN"
)

const (
	DeploymentConfigJobArgName  = "deployment_config"
	CurrentContainersJobArgName = "current_containers"
	NewContainersJobArgName     = "new_containers"
)

func createDeploymentJob(config config.DeploymentConfig, action action) (job.Job, error) {
	switch action {
	case Up:
		return createContainersJob(config), nil
	case Down:
		return deleteContainersJob(config), nil
	default:
		return job.Job{}, fmt.Errorf("unknown deployment action %s", action)
	}
}

func createContainersJob(config config.DeploymentConfig) job.Job {
	currContainers := make([]container.KraneContainer, 0)
	newContainers := make([]container.KraneContainer, 0)
	retryPolicy := utils.UIntEnv("DEPLOYMENT_RETRY_POLICY")

	jobsArgs := job.Args{
		DeploymentConfigJobArgName:  config,
		CurrentContainersJobArgName: &currContainers,
		NewContainersJobArgName:     &newContainers,
	}

	return job.Job{
		ID:          uuid.Generate().String(),
		Namespace:   config.Name,
		Type:        ContainerCreate,
		Args:        jobsArgs,
		RetryPolicy: retryPolicy,
		Run:         createContainerResources,
	}
}

func deleteContainersJob(config config.DeploymentConfig) job.Job {
	containers := make([]container.KraneContainer, 0)
	retryPolicy := utils.UIntEnv("DEPLOYMENT_RETRY_POLICY")

	jobsArgs := job.Args{
		DeploymentConfigJobArgName:  config,
		CurrentContainersJobArgName: &containers,
	}

	return job.Job{
		ID:          uuid.Generate().String(),
		Namespace:   config.Name,
		Type:        ContainerDelete,
		Args:        jobsArgs,
		RetryPolicy: retryPolicy,
		Run:         deleteContainerResources,
	}
}
