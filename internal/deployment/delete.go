package deployment

import (
	"context"

	"github.com/biensupernice/krane/docker"

	"github.com/biensupernice/krane/internal/logger"
)

// Remove : a deployments resources
func Remove(ctx *context.Context, t Template) {
	go EmitEvent("Removing deployment resources", t)

	retries := 3
	for i := 0; i < retries; i++ {
		logger.Debugf("Attempt [%d] to remove %s resources", i+1, t.Name)

		err := deleteDockerResources(ctx, t)
		if err != nil {
			logger.Debugf("Unable to remove resouces for %s - %s", t.Name, err.Error())
			logger.Debug("Waiting 10 seconds before retrying")
			wait(10)
			continue
		}
		break
	}

	go EmitEvent("Finished removing resources", t)
	logger.Debugf("Finished removing resources for - %s", t.Name)
}

func deleteDockerResources(ctx *context.Context, t Template) (err error) {
	containers := GetContainers(ctx, t.Name)
	for _, container := range containers {
		// Stop the container
		err = docker.StopContainer(ctx, container.ID)
		if err != nil {
			logger.Debugf("Unable to stop %s - %s", container.ID, err.Error())
			return
		}
		logger.Debugf("Stopped container %s", container.ID)

		// Remove the container
		err = docker.RemoveContainer(ctx, container.ID)
		if err != nil {
			logger.Debugf("Unable to remove %s - %s", container.ID, err.Error())
			return
		}
		logger.Debugf("Removed container %s", container.ID)

		// Delete docker image
		_, err = docker.RemoveImage(ctx, container.ImageID)
		if err != nil {
			logger.Debugf("Unable to remove image %s - %s", container.Image, err.Error())
			return
		}
		logger.Debugf("Removed image %s", container.Image)
	}

	logger.Debugf("Cleaned up docker resources for %s", t.Name)
	return
}
