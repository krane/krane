package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/biensupernice/krane/data"

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

// Deployment :
type Deployment struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`                     // Deployment name
	Registry      string   `json:"registry"`                 // Docker registry url
	Image         string   `json:"image" binding:"required"` // Docker image name
	Tag           string   `json:"tag"`                      // Docker image tag
	ContainerPort string   `json:"container_port"`           // Port to expose from the container
	HostPort      string   `json:"host_port"`                // Port to expose to the host
	Metadata      Metadata `json:"metadata"`                 // Metadata about the deployment
}

// Metadata : about a deployment
type Metadata struct {
	Status      string  `json:"status"`       // Deployment status
	ContainerID string  `json:"container_id"` // Deployment container id
	Events      []Event `json:"events"`       // Deployment info, error
}

// Event : specify a deployment event
type Event struct {
	Timestamp time.Time         `json:"timestamp" binding:"required"`
	Data      map[string]string `json:"data"`
}

// StartDeployment : start a deployment
func StartDeployment(d *Deployment) {
	// Set deployments defaults
	setDefaults(d)

	// Start deployment process
	updateDeploymentStatus(d.ID, StatusInProgress)
	addDeploymentEvent(d.ID, &Event{Timestamp: time.Now(), Data: map[string]string{"message": fmt.Sprintf("Starting Deployment - %s", d.ID)}})

	// Store deployment
	err := storeDeployment(d)
	if err != nil {
		log.Println(err.Error())
		addDeploymentEvent(d.ID, &Event{Timestamp: time.Now(), Data: map[string]string{"message": fmt.Sprintf("Unable to save deployment")}})
		return
	}

	// This number represent the amount of tries krane will attempt
	// to start a deployment before marking it as failed
	attemptsBeforeFailing := 3
	for attempts := 0; attempts < attemptsBeforeFailing; attempts++ {
		_, err := Deploy(d)
		if err != nil {
			// Deployment failed, update status, record event
			updateDeploymentStatus(d.ID, StatusFailed)
			addDeploymentEvent(d.ID, &Event{
				Timestamp: time.Now(),
				Data:      map[string]string{"message": fmt.Sprintf("[%d/%d] Deployment failed - %s", attempts+1, attemptsBeforeFailing, err.Error())}})

			// Return if retry limit has exceeded
			if attempts == attemptsBeforeFailing-1 {
				addDeploymentEvent(d.ID, &Event{
					Timestamp: time.Now(),
					Data:      map[string]string{"message": fmt.Sprintf("Exceeded retry limit of %d, stopping deployment", attemptsBeforeFailing)}})
				return
			}

			time.Sleep(10 * time.Second) // 10 seconds
			continue
		}

		updateMetadata(d.ID, d.Metadata)
		addDeploymentEvent(d.ID, &Event{Timestamp: time.Now(), Data: map[string]string{"message": fmt.Sprintf(fmt.Sprintf("Deployment finished [Container ID] %s", d.Metadata.ContainerID))}})
		break
	}
}

func updateMetadata(id string, m Metadata) {
	// Get deployment from db
	d := *getDeployment(id)

	// Update deployment Status
	d.Metadata.ContainerID = m.ContainerID
	d.Metadata.Status = m.Status
	d.Metadata.Events = append(m.Events)

	// Update deployment in ddb
	dBytes, _ := json.Marshal(d)
	data.Put(data.DeploymentsBucket, id, dBytes)
}

func getDeployment(id string) *Deployment {
	var d Deployment
	dBytes := data.Get(data.DeploymentsBucket, id)
	json.Unmarshal(dBytes, &d)
	return &d
}

func storeDeployment(d *Deployment) error {
	// Store deployment
	dplmntBytes, _ := json.Marshal(d)
	return data.Put(data.DeploymentsBucket, d.ID, dplmntBytes)
}

func addDeploymentEvent(id string, event *Event) error {
	// Get deployment from db
	d := *getDeployment(id)

	// Update deployment Status
	d.Metadata.Events = append(d.Metadata.Events, *event)

	// Update deployment in ddb
	dBytes, _ := json.Marshal(d)
	return data.Put(data.DeploymentsBucket, id, dBytes)
}

func updateDeploymentStatus(id, status string) {
	// Get deployment from db
	d := *getDeployment(id)

	// Update deployment Status
	d.Metadata.Status = status

	// Update deployment in ddb
	dBytes, _ := json.Marshal(d)
	data.Put(data.DeploymentsBucket, id, dBytes)
}

func updateDeploymentContainerID(id, containerID string) {
	// Get deployment from db
	d := *getDeployment(id)

	// Update deployment Status
	d.Metadata.ContainerID = containerID

	// Update deployment in ddb
	dBytes, _ := json.Marshal(d)
	data.Put(data.DeploymentsBucket, id, dBytes)
}

func setDefaults(deployment *Deployment) *Deployment {
	// Set deployment uid if not set
	id, _ := uuid.NewUUID()
	deployment.ID = id.String()

	// Set deployment name to {image}-{id} name if not provided
	if deployment.Name == "" {
		deployment.Name = fmt.Sprintf("%s-%s", "bsn", deployment.ID)
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
func Deploy(deployment *Deployment) (*Deployment, error) {
	metadata := &deployment.Metadata
	metadata.Status = StatusInProgress // Deployment : InProgress
	updateDeploymentStatus(deployment.ID, metadata.Status)

	// Get docker client
	_, err := New()
	if err != nil {
		metadata.Status = StatusFailed
		addDeploymentEvent(deployment.ID, &Event{
			Timestamp: time.Now(),
			Data:      map[string]string{"message": err.Error()}})
		return deployment, err
	}

	// Start deployment context
	ctx := context.Background()

	// Format docker image url source
	img := FormatImageSourceURL(deployment.Registry, deployment.Image, deployment.Tag)
	addDeploymentEvent(deployment.ID, &Event{
		Timestamp: time.Now(),
		Data:      map[string]string{"message": fmt.Sprintf("Pulling image: %s", img)}})

	// Pull docker image
	err = PullImage(&ctx, img)
	if err != nil {
		metadata.Status = StatusFailed
		return deployment, err
	}

	// Create docker container
	createContainerResp, err := CreateContainer(&ctx,
		img,
		deployment.Name,
		deployment.HostPort,
		deployment.ContainerPort)
	if err != nil {
		metadata.Status = StatusFailed
		return deployment, err
	}

	// Set deployment metadata container id
	metadata.ContainerID = createContainerResp.ID

	// Start docker container
	err = StartContainer(&ctx, metadata.ContainerID)
	if err != nil {
		RemoveContainer(&ctx, metadata.ContainerID)
		metadata.Status = StatusFailed
		return deployment, err
	}

	metadata.Status = StatusSucceeded
	addDeploymentEvent(deployment.ID, &Event{
		Timestamp: time.Now(),
		Data:      map[string]string{"message": fmt.Sprintf("Succesfully deployed %s - %s", deployment.Name, metadata.ContainerID)}})

	return deployment, nil
}
