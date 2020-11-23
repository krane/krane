package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/secrets"
)

func createContainerResources(args job.Args) error {
	wf := newWorkflow("CreateContainerResources", args)

	wf.with("GetCurrentContainers", getCurrentContainers)
	wf.with("CreateSecretsCollection", createSecretsCollection)
	wf.with("CreateJobsCollection", createJobsCollection)
	wf.with("PullImage", pullImage)
	wf.with("CreateContainers", createContainers)
	wf.with("StartContainers", startContainers)
	wf.with("CheckNewContainersHealth", checkNewContainersHealth)
	wf.with("RemoveOldContainers", cleanupCurrentContainers)

	return wf.start()
}

func pullImage(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)

	ctx := context.Background()
	defer ctx.Done()

	return docker.GetClient().PullImage(ctx, cfg.Registry, cfg.Image, cfg.Tag)
}

func createContainers(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)

	newContainers := make([]container.Kcontainer, 0)
	for i := 0; i < cfg.Scale; i++ {
		newContainer, err := container.Create(cfg)
		if err != nil {
			return err
		}
		newContainers = append(newContainers, newContainer)
	}
	logrus.Debugf("Created %d containers", len(newContainers))
	args["newContainers"] = &newContainers
	return nil
}

func startContainers(args job.Args) error {
	newContainers := args["newContainers"].(*[]container.Kcontainer)
	containersStarted := 0
	for _, newContainer := range *newContainers {
		err := newContainer.Start()
		if err != nil {
			return err
		}
		containersStarted++
	}
	logrus.Debugf("Started %d containers", containersStarted)
	return nil
}

func checkNewContainersHealth(args job.Args) error {
	newContainers := args["newContainers"].(*[]container.Kcontainer)

	pollRetry := 10
	for _, newContainer := range *newContainers {
		for i := 0; i <= pollRetry; i++ {
			expBackOff := time.Duration(10 * i)
			time.Sleep(expBackOff * time.Second)

			ok, err := newContainer.Ok()
			if err != nil {
				if i == pollRetry {
					return fmt.Errorf("container health unstable %v", err)
				}
				continue
			}

			if !ok {
				if i == pollRetry {
					return fmt.Errorf("container health unstable %v", err)
				}
				continue
			}

			// If here container health should be healthy
			break
		}
	}

	return nil
}

func createSecretsCollection(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)
	return secrets.CreateCollection(cfg.Name)
}

func createJobsCollection(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)
	return job.CreateCollection(cfg.Name)
}
