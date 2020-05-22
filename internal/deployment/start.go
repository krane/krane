package deployment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/biensupernice/krane/docker"

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
		err := deployWithDocker(&t)
		if err != nil {
			log.Printf("Unable to tsart deployment %s", err.Error())

			wait(10)
			continue
		}
	}
}

// deployWithDocker : workflow to deploy a docker container
func deployWithDocker(t *Template) (err error) {
	// deployment context
	ctx := context.Background()

	// create well formated url to fetch docker image
	img := docker.FormatImageSourceURL(t.Config.Registry, t.Config.Image, t.Config.Tag)

	// Pull docker image
	err = docker.PullImage(&ctx, img)
	if err != nil {
		return err
	}

	// Create docker container
	dID := uuid.NewSHA1(uuid.New(), []byte(t.Name)) // deployment ID
	containerName := fmt.Sprintf("%s-%s", t.Name, dID)
	createContainerResp, err := docker.CreateContainer(
		&ctx,
		img,
		containerName,
		t.Config.HostPort,
		t.Config.ContainerPort)
	if err != nil {
		return
	}

	containerID := createContainerResp.ID

	// Start docker container
	err = docker.StartContainer(&ctx, containerID)
	if err != nil {
		docker.RemoveContainer(&ctx, containerID)
		return
	}

	return
}

func wait(s time.Duration) {
	time.Sleep(s * time.Second)
}
