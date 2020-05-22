package deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/logger"

	"github.com/google/uuid"
)

const (
	// StatusFailed : deployment failed
	StatusFailed = "Failed"

	// StatusSucceeded : deployment succeeded
	StatusSucceeded = "Succeeded"

	// StatusInProgress : deployment in progress
	StatusInProgress = "InProgress"
)

// Start : a deployment using a template
func Start(t Template) {
	retries := 3
	for i := 0; i < retries; i++ {
		logger.Debugf("Attempt [%d] to deploy %s", i+1, t.Name)
		err := deployWithDocker(&t)
		if err != nil {
			logger.Debugf("Unable to start deployment %s", err.Error())
			logger.Debug("Waiting 10 seconds before retrying")
			wait(10)
			continue
		}
		break
	}
	logger.Debugf("Deployment complete - %s", t.Name)
}

// deployWithDocker : workflow to deploy a docker container
func deployWithDocker(t *Template) (err error) {
	// deployment context
	ctx := context.Background()

	// create well formated url to fetch docker image
	img := docker.FormatImageSourceURL(t.Config.Registry, t.Config.Image, t.Config.Tag)
	logger.Debugf("Puling %s", img)

	// Pull docker image
	err = docker.PullImage(&ctx, img)
	if err != nil {
		logger.Debugf("Unable to pull the image - %s", err.Error())
		return err
	}

	// Create docker container
	dID := uuid.NewSHA1(uuid.New(), []byte(t.Name)) // deployment ID
	shortID := dID.String()[0:8]
	containerName := fmt.Sprintf("%s-%s", t.Name, shortID)
	createContainerResp, err := docker.CreateContainer(
		&ctx,
		img,
		containerName,
		t.Config.HostPort,
		t.Config.ContainerPort)
	if err != nil {
		logger.Debugf("Unable to create docker container - %s", err.Error())
		return
	}

	containerID := createContainerResp.ID
	logger.Debugf("Container created with id %s", containerName)

	// Start docker container
	err = docker.StartContainer(&ctx, containerID)
	if err != nil {
		logger.Debugf("Unable to start container - %s", err.Error())
		docker.RemoveContainer(&ctx, containerID)
		return
	}
	logger.Debugf("Container started with the name %s", containerName)
	return
}

func wait(s time.Duration) {
	time.Sleep(s * time.Second)
}
