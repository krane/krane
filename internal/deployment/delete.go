package deployment

import (
	"context"
	"fmt"

	"github.com/biensupernice/krane/docker"

	"github.com/biensupernice/krane/internal/container"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/spec"
)

// Remove : a deployments resources
func Remove(ctx *context.Context, s spec.Spec) (success bool) {
	go EmitEvent("Removing deployment resources", s)

	status := DeletingStatus
	retries := 5
	for i := 0; i < retries; i++ {
		logger.Debugf("Attempt [%d] to remove %s resources", i+1, s.Name)

		err := deleteDockerResources(ctx, s)
		if err != nil {
			status = FailedStatus

			logger.Debugf("Unable to remove resouces for %s - %s", s.Name, err.Error())
			logger.Debug("Waiting 10 seconds before retrying")
			wait(10)
			continue
		}

		status = ReadyStatus
		break
	}

	// If deployment errored out log failure event and success false
	if status != ReadyStatus {
		errMsg := fmt.Sprintf("Unable to remove deployment %s", s.Name)
		go EmitEvent(errMsg, s)
		logger.Debugf(errMsg)
		return false
	}

	// If deployment resources succesfully got removed, log event and return true
	successMsg := fmt.Sprintf("Succesfully deleted deployment %s", s.Name)
	logger.Debugf(successMsg)
	go EmitEvent(successMsg, s)
	return true
}

func deleteDockerResources(ctx *context.Context, s spec.Spec) (err error) {
	containers := container.Get(ctx, s.Name)
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

	logger.Debugf("Cleaned up docker resources for %s", s.Name)
	return
}
