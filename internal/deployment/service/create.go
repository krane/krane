package service

import (
	"context"
	"errors"
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
	cfg := args["config"].(config.Config)
	logrus.Debugf("Starting deployment workflow for %s", cfg.Name)

	client := docker.GetClient()

	// 1. get curr containers
	currContainer, err := container.GetKontainersByNamespace(client, cfg.Name)
	if err != nil {
		return err
	}
	logrus.Debugf("Found %d existing containers for %s", len(currContainer), cfg.Name)

	// 2. pull image
	image := docker.FormatImageSourceURL(cfg.Registry, cfg.Image, cfg.Tag)
	if err := pullImage(image); err != nil {
		return err
	}
	logrus.Debugf("Pulled %s for %s", image, cfg.Name)

	// 3. create containers
	kontainer, createContainerErr := container.Create(cfg)
	if createContainerErr != nil {
		return createContainerErr
	}
	logrus.Debugf("Created containers for %s", cfg.Name)

	// 4. start containers
	startContainerErr := kontainer.Start()
	if startContainerErr != nil {
		return startContainerErr
	}
	logrus.Debugf("Started containers for %s", cfg.Name)

	// 5. check container health
	containerHealthCheckErr := pollContainerUntilHealthy(kontainer)
	if containerHealthCheckErr != nil {
		return err
	}
	logrus.Debugf("Healthy containers for %s", cfg.Name)

	// 6. cleanup old containers
	for _, c := range currContainer {
		if err := c.Remove(); err != nil {
			return err
		}
	}
	logrus.Debugf("Removed %d containers for %s", len(currContainer), cfg.Name)

	// 7. cleanup old images
	// TODO:

	logrus.Debugf("Completed deployment workflow for %s", cfg.Name)
	return nil
}

func pullImage(image string) error {
	ctx := context.Background()
	client := docker.GetClient()
	err := client.PullImage(&ctx, image)
	ctx.Done()
	return err
}

func pollContainerUntilHealthy(container container.Kontainer) error {
	pollRetry := 10
	for i := 0; i <= pollRetry; i++ {
		expBackOff := time.Duration(10 * i)
		time.Sleep(expBackOff * time.Second)

		ok, err := container.Ok()
		if err != nil {
			continue
		}

		if !ok {
			continue
		}

		return nil
	}

	return errors.New("unable to determine container health")
}
