package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/biensupernice/krane/data"

	"github.com/google/uuid"
)

const (
	// DeploymentStatusFailed : deployment failed
	DeploymentStatusFailed = "Failed"

	// DeploymentStatusSucceeded : deployment succeeded
	DeploymentStatusSucceeded = "Succeeded"

	// DeploymentStatusPending : deployment queued but not started.
	DeploymentStatusPending = "Pending"

	// DeploymentStatusInProgress : deployment in progress
	DeploymentStatusInProgress = "InProgress"
)

// Deployment : state of a deployment
type Deployment struct {
	ID            string `json:"id"`
	Name          string `json:"name"`                     // Deployment name
	Status        string `json:"status"`                   // Deployment status
	ContainerID   string `json:"container_id"`             // Deployment container id
	Registry      string `json:"registry"`                 // Docker registry url
	Image         string `json:"image" binding:"required"` // Docker image name
	Tag           string `json:"tag"`                      // Docker image tag
	ContainerPort string `json:"container_port"`           // Port to expose from the container
	HostPort      string `json:"host_port"`                // Port to expose to the host
}

// QueueDeployment : queue a deployment
func QueueDeployment(deployment Deployment) {
	// Set defaults
	dplmnt := setDeploymentDefaults(&deployment)

	log.Printf("Queuing deployment - %s", dplmnt.Name)

	// Set deployment status to pending
	dplmnt.Status = DeploymentStatusPending

	// Insert the current deployment into the deployments buckets
	dplmntBytes, _ := json.Marshal(&dplmnt)
	data.Put(data.DeploymentsBucket, dplmnt.ID, dplmntBytes)

	// TODO: Queue up deployment
}

func setDeploymentDefaults(deployment *Deployment) *Deployment {
	// Set deployment uid if not set
	id, _ := uuid.NewUUID()
	deployment.ID = id.String()

	// Set deployment name to {image}-{id} name if not provided
	if deployment.Name == "" {
		deployment.Name = fmt.Sprintf("%s-%s", deployment.Image, deployment.ID)
	}

	// Set docker registry if not provided
	if deployment.Registry == "" {
		deployment.Registry = "registry.hub.docker.com"
	}

	// Set image tag to `latest` if not provided
	if deployment.Tag == "" {
		deployment.Tag = "latest"
	}

	return deployment
}

// Deploy : docker container
func Deploy(deployment Deployment) (containerID string, err error) {
	deployment.Status = DeploymentStatusInProgress
	log.Printf("Deploying %s\n", deployment.Name)

	// Get docker client
	_, err = New()
	if err != nil {
		return
	}

	// Start deployment context
	ctx := context.Background()

	// Format docker image url source
	img := FormatImageSourceURL(deployment.Registry, deployment.Image, deployment.Tag)
	log.Printf("Pulling image: %s\n", img)

	// Pull docker image
	err = PullImage(&ctx, img)
	if err != nil {
		return
	}

	// Create docker container
	createContainerResp, err := CreateContainer(&ctx,
		img,
		deployment.Name,
		deployment.HostPort,
		deployment.ContainerPort)
	containerID = createContainerResp.ID
	if err != nil {
		return
	}

	// Docker start container
	err = StartContainer(&ctx, containerID)
	if err != nil {
		// If error starting the container, remove it
		RemoveContainer(&ctx, containerID)
		return
	}

	deployment.Status = DeploymentStatusSucceeded
	log.Printf("Deployed %s - 📦 %s\n", deployment.Name, containerID)

	return
}
