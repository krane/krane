package deployment

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/spec"
	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
)

// Deployment :
type Deployment struct {
	Spec       spec.Spec         `json:"spec" binding:"required"`
	Containers []types.Container `json:"containers"`
}

// Deployment statuseseseses
const (
	// ReadyStatus : Deployment is ready
	ReadyStatus = "Ready"
	// InProgressStatus : Deployment is in progress
	InProgressStatus = "InProgress"
	// FailedStatus : Deployment has failed
	FailedStatus = "Failed"
	// DeletingStatus : Deployment is being deleting along with its resources
	DeletingStatus = "Deleting"
)

// Start : a deployment using a template and the tag that will be used for
// the image that will deployed
func Start(ctx *context.Context, s spec.Spec, tag string) {
	go EmitEvent("Starting deployment", s)

	status := InProgressStatus
	retries := 3
	for i := 0; i < retries; i++ {
		logger.Debugf("Attempt [%d] to create %s resources", i+1, s.Name)

		containerID, err := deployWithDocker(ctx, s, tag)
		if err != nil {
			status = FailedStatus

			logger.Debugf("Unable to start deployment %s", err.Error())
			logger.Debug("Waiting 10 seconds before retrying")
			wait(10)
			continue
		}

		// Check if container ID was returned
		if containerID == "" {
			status = FailedStatus

			logger.Debug("containerID not returned from deployment attempt, retrying")
			wait(10)
			continue
		}

		status = ReadyStatus
		break
	}

	msg := fmt.Sprintf("Deployment %s", status)
	go EmitEvent(msg, s)
	logger.Debugf(msg)
}

// deployWithDocker : workflow to deploy a docker container
func deployWithDocker(ctx *context.Context, s spec.Spec, tag string) (containerID string, err error) {
	// create well formated url to fetch docker image
	img := docker.FormatImageSourceURL(s.Config.Registry, s.Config.Image, tag)
	logger.Debugf("Pulling %s", img)

	// Pull docker image
	go EmitEvent("Pulling image", s)
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
	go EmitEvent("Creating the container", s)
	dID := uuid.NewSHA1(uuid.New(), []byte(s.Name)) // deployment ID
	shortID := dID.String()[0:8]
	containerName := fmt.Sprintf("%s-%s", s.Name, shortID)
	createContainerResp, err := docker.CreateContainer(
		ctx,
		img,
		s.Name,
		containerName,
		netID,
		s.Config.HostPort,
		s.Config.ContainerPort)
	if err != nil {
		logger.Debugf("Unable to create docker container - %s", err.Error())
		return
	}

	containerID = createContainerResp.ID
	logger.Debugf("Container created with id %s", containerID)

	// Connect container to network
	err = docker.ConnectContainerToNetwork(ctx, netID, containerID)
	if err != nil {
		logger.Debugf("Unable to connect container to docker network - %s", err.Error())
		return
	}

	// Start docker container
	go EmitEvent("Starting container", s)
	err = docker.StartContainer(ctx, containerID)
	if err != nil {
		logger.Debugf("Unable to start container - %s", err.Error())
		docker.RemoveContainer(ctx, containerID)
		return
	}

	go EmitEvent("Container started", s)
	logger.Debugf("Container started with the name %s", containerName)

	return
}

func wait(s time.Duration) {
	time.Sleep(s * time.Second)
}
