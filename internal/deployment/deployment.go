package deployment

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/event"
	"github.com/biensupernice/krane/internal/deployment/spec"
	"github.com/biensupernice/krane/logger"
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
	go event.Emit("Starting deployment", s)

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
	go event.Emit(msg, s)
	logger.Debugf(msg)
}

// Remove : a deployments resources
func Remove(ctx *context.Context, s spec.Spec) {
	go event.Emit("Removing deployment resources", s)

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

	// If deployment errored out log failure event
	if status != ReadyStatus {
		errMsg := fmt.Sprintf("Unable to remove deployment %s", s.Name)
		go event.Emit(errMsg, s)
		logger.Debugf(errMsg)
		return
	}

	// If deployment resources succesfully got removed, log event
	successMsg := fmt.Sprintf("Succesfully deleted deployment %s", s.Name)
	logger.Debugf(successMsg)
	go event.Emit(successMsg, s)

	// Delete deployment spec ONLY if succesfully removed all deployment resources
	s.Delete()

	return
}

// deployWithDocker : workflow to deploy a docker container
func deployWithDocker(ctx *context.Context, s spec.Spec, tag string) (containerID string, err error) {
	// create well formated url to fetch docker image
	img := docker.FormatImageSourceURL(s.Config.Registry, s.Config.Image, tag)
	logger.Debugf("Pulling %s", img)

	// Pull docker image
	go event.Emit("Pulling image", s)
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
	go event.Emit("Creating the container", s)
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
	go event.Emit("Starting container", s)
	err = docker.StartContainer(ctx, containerID)
	if err != nil {
		logger.Debugf("Unable to start container - %s", err.Error())
		docker.RemoveContainer(ctx, containerID)
		return
	}

	go event.Emit("Container started", s)
	logger.Debugf("Container started with the name %s", containerName)

	return
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

// Helper to wait an X amount of seconds
func wait(s time.Duration) { time.Sleep(s * time.Second) }
