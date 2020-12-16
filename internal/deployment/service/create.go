package service

import (
	"context"
	"fmt"
	"time"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/secrets"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
)

// createContainerResources : creates a workflow that defines the container creation process
func createContainerResources(args job.Args) error {
	wf := job.NewWorkflow("CreateContainerResources", args)

	wf.With("GetCurrentContainers", getCurrentContainers)
	wf.With("EnsureSecretsCollection", ensureSecretsCollection)
	wf.With("EnsureJobsCollection", ensureJobsCollection)
	wf.With("PullDockerImage", pullImage)
	wf.With("CreateContainers", createContainers)
	wf.With("StartContainers", startContainers)
	wf.With("PerformContainerHealth", checkNewContainersHealth)
	wf.With("CleanupOldContainers", removeCurrentContainers)

	return wf.Start()
}

func pullImage(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)

	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().PullImage(ctx, cfg.Registry, cfg.Image, cfg.Tag)
}

func createContainers(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)

	containers := make([]container.KraneContainer, 0)
	for i := 0; i < cfg.Scale; i++ {
		c, err := container.Create(cfg)
		if err != nil {
			return err
		}
		containers = append(containers, c)
	}

	args[NewContainersJobArgName] = &containers
	logger.Debugf("Created %d container(s)", len(containers))
	return nil
}

func startContainers(args job.Args) error {
	containers := args.GetArg(NewContainersJobArgName).(*[]container.KraneContainer)
	count := 0
	for _, c := range *containers {
		if err := c.Start(); err != nil {
			return err
		}
		count++
	}
	logger.Debugf("Started %d container(s)", count)
	return nil
}

func checkNewContainersHealth(args job.Args) error {
	containers := args.GetArg(NewContainersJobArgName).(*[]container.KraneContainer)

	pollRetry := 10
	for _, c := range *containers {
		for i := 0; i <= pollRetry; i++ {
			expBackOff := time.Duration(10 * i)
			time.Sleep(expBackOff * time.Second)

			ok, err := c.Ok()
			if err != nil {
				if i == pollRetry {
					return fmt.Errorf("container is not healthy %v", err)
				}
				continue
			}

			if !ok {
				if i == pollRetry {
					return fmt.Errorf("container is not healthy %v", err)
				}
				continue
			}

			// if reached here, container healthy
			break
		}
	}
	return nil
}

func ensureSecretsCollection(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)
	return secrets.CreateCollection(cfg.Name)
}

func ensureJobsCollection(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)
	return job.CreateCollection(cfg.Name)
}
