package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
)

// Run: The handler in charge of creating new deployment container resources
//
// Args: Krane Config
//
// Signature: func Run(args Args) error
//
// 1. Fetch current container for the deployment. The deployment metadata is provided in the args.
//
// 2. Create new container(s)
// - Every container is created suffixed with a unique id, format {name}-{shortuuid} for example api-Cekw67uyMpBGZLRP2HFVbe
//
// 3. Up new container(s)
// - New containers are automatically attached to the traefick network
// - Since cleanup has not occurred, traefik is load balancing between old container and new created containers
//
// 4. Poll new containers until (a) healthy | (b) unhealthy
// a. if healthy, create a new job to deleteContainerResources previous containers if > 0
// b. if unhealthy or state "unknown" for X amount of time, remove newly created containers and report failure.
// - Note: when new containers are in "unhealthy" state the previous containers are not cleaned up unless
// - previous containers are also unhealthy

// createContainerResources create containers deployment workfllow
func createContainerResources(args job.Args) error {
	konfig := args["config"].(config.Config)
	logrus.Debugf("Starting deployment workflow for %s", konfig.Name)

	client := docker.GetClient()

	// 1. get curr containers
	currContainer, err := container.GetKontainersByNamespace(client, konfig.Name)
	if err != nil {
		return err
	}
	logrus.Debugf("Found %d existing containers for %s", len(currContainer), konfig.Name)

	// 2. pull image
	image := docker.FormatImageSourceURL(konfig.Registry, konfig.Image, konfig.Tag)
	pullImageErr := pullImage(image)
	if pullImageErr != nil {
		return pullImageErr
	}

	// 3. create containers
	createContainerErr := createContainers()
	if createContainerErr != nil {
		return createContainerErr
	}
	// 4. createContainerResources containers
	startContainerErr := startContainers()
	if startContainerErr != nil {
		return startContainerErr
	}

	// 5. check container health

	// 6. cleanup

	time.Sleep(30 * time.Second)
	return nil
}

func createContainers() error {
	return nil
}

func startContainers() error {
	return nil
}

func pullImage(image string) error {
	ctx := context.Background()
	client := docker.GetClient()
	err := client.PullImage(&ctx, image)
	ctx.Done()
	return err
}
