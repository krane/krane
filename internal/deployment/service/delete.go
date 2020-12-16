package service

import (
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/secrets"
	"github.com/biensupernice/krane/internal/job"
)

func deleteContainerResources(args job.Args) error {
	wf := job.NewWorkflow("DeleteContainerResources", args)

	wf.With("GetCurrentContainers", getCurrentContainers)
	wf.With("RemoveContainers", removeCurrentContainers)
	wf.With("RemoveDeploymentSecrets", deleteDeploymentSecrets)
	wf.With("RemoveDeploymentJobs", deleteDeploymentJobs)
	wf.With("RemoveDeploymentConfig", deleteDeploymentConfig)

	return wf.Start()
}

func removeCurrentContainers(args job.Args) error {
	containers := args.GetArg(CurrentContainersJobArgName).(*[]container.KraneContainer)
	for _, c := range *containers {
		err := c.Remove()
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteDeploymentSecrets(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)
	return secrets.DeleteCollection(cfg.Name)
}

func deleteDeploymentJobs(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)
	return job.DeleteCollection(cfg.Name)
}

func deleteDeploymentConfig(args job.Args) error {
	cfg := args.GetArg(DeploymentConfigJobArgName).(config.DeploymentConfig)
	return config.Delete(cfg.Name)
}
