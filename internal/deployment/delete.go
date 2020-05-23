package deployment

import (
	"context"
	"strings"

	"github.com/biensupernice/krane/docker"

	"github.com/biensupernice/krane/internal/logger"
)

// DeleteDockerResources : delete a deployments docker resources
func DeleteDockerResources(ctx *context.Context, t Template) {
	retries := 3
	for i := 0; i < retries; i++ {
		logger.Debugf("Attempt [%d] to delete deployment %s", i, t.Name)

		// Delete docker resources
		err := deleteDockerResources(ctx, t)
		if err != nil {
			logger.Debugf("Unable to remove containers for deployment %s - %s", t.Name, err.Error())
			continue
		}

		logger.Debugf("Removed container resources for deployment %s", t.Name)
		break
	}
}

func deleteDockerResources(ctx *context.Context, t Template) (err error) {
	// Get all containers
	containers, err := docker.ListContainers(ctx)
	if err != nil {
		return err
	}

	logger.Debugf("Received %d containers", len(containers))

	// Find out which containers belong to the deployment
	for i := 0; i < len(containers); i++ {
		c := containers[i]

		// Check the containers names and
		// if the prefix matches with the deployment name
		// remove that container
		for n := 0; n < len(c.Names); n++ {
			name := c.Names[n]
			if strings.HasPrefix(name[1:len(name)], t.Name) {
				// Stop the container
				err = docker.StopContainer(ctx, containers[i].ID)
				if err != nil {
					logger.Debugf("Unable to stop %s - %s", name, err.Error())
					return err
				}

				// Remove the container
				err = docker.RemoveContainer(ctx, containers[i].ID)
				if err != nil {
					logger.Debugf("Unable to remove %s - %s", name, err.Error())
					return err
				}

				// Delete docker image
				_, err = docker.RemoveImage(ctx, c.ImageID)
				if err != nil {
					logger.Debugf("Unable to remove image %s - %s", c.Image, err.Error())
					return err
				}

				logger.Debugf("Cleaned up docker resources for %s", name)
			}
		}
	}
	return
}
