package service

import (
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/job"
)

func getCurrentContainers(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)

	containers, err := container.GetKraneContainersByDeployment(cfg.Name)
	if err != nil {
		return err
	}

	args[CurrentContainersJobArgName] = &containers
	return nil
}
