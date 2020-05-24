package deployment

import (
	"context"
	"errors"

	"fmt"
	"os"
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

// Start : a deployment using a template and the tag that will be used for
// the image that will deployed
func Start(ctx *context.Context, t Template, tag string) {
	retries := 3
	for i := 0; i < retries; i++ {
		logger.Debugf("Attempt [%d] to deploy %s", i+1, t.Name)

		containerID, err := deployWithDocker(ctx, t, tag)
		if err != nil {
			logger.Debugf("Unable to start deployment %s", err.Error())
			logger.Debug("Waiting 10 seconds before retrying")
			wait(10)
			continue
		}

		// Check if container ID was returned
		if containerID == "" {
			logger.Debug("containerID not returned from deployment attempt, retrying")
			wait(10)
			continue
		}
		break
	}

	logger.Debugf("Deployment complete - %s", t.Name)
}

// deployWithDocker : workflow to deploy a docker container
func deployWithDocker(ctx *context.Context, t Template, tag string) (containerID string, err error) {
	// create well formated url to fetch docker image
	img := docker.FormatImageSourceURL(t.Config.Registry, t.Config.Image, tag)
	logger.Debugf("Pulling %s", img)

	// Pull docker image
	err = docker.PullImage(ctx, img)
	if err != nil {
		logger.Debugf("Unable to pull the image - %s", err.Error())
		return
	}

	// Krane Network ID to connect the container
	netID := os.Getenv("KRANE_NETWORK_ID")
	if netID == "" {
		return "", errors.New("Unable to create docker container, krane network not found")
	}

	// Create docker container
	dID := uuid.NewSHA1(uuid.New(), []byte(t.Name)) // deployment ID
	shortID := dID.String()[0:8]
	containerName := fmt.Sprintf("%s-%s", t.Name, shortID)
	createContainerResp, err := docker.CreateContainer(
		ctx,
		img,
		containerName,
		netID,
		t.Config.HostPort,
		t.Config.ContainerPort)
	if err != nil {
		logger.Debugf("Unable to create docker container - %s", err.Error())
		return
	}

	containerID = createContainerResp.ID
	logger.Debugf("Container created with id %s", containerID)

	// Start docker container
	err = docker.StartContainer(ctx, containerID)
	if err != nil {
		logger.Debugf("Unable to start container - %s", err.Error())
		docker.RemoveContainer(ctx, containerID)
		return
	}
	logger.Debugf("Container started with the name %s", containerName)

	return
}

func wait(s time.Duration) {
	time.Sleep(s * time.Second)
}
