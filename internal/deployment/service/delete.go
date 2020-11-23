package service

import (
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/secrets"
)

func deleteContainerResources(args job.Args) error {
	wf := newWorkflow("DeleteContainerResources", args)

	wf.with("GetCurrentContainers", getCurrentContainers)
	wf.with("RemoveContainers", cleanupCurrentContainers)
	wf.with("RemoveDeploymentSecrets", deleteDeploymentSecrets)
	wf.with("RemoveDeploymentJobs", deleteDeploymentJobs)
	wf.with("RemoveDeploymentConfig", deleteDeploymentConfig)

	return wf.start()
}

func cleanupCurrentContainers(args job.Args) error {
	currContainers := args["currContainers"].(*[]container.Kcontainer)
	for _, oldContainer := range *currContainers {
		err := oldContainer.Remove()
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteDeploymentSecrets(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)
	return secrets.DeleteCollection(cfg.Name)
}

func deleteDeploymentJobs(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)
	return job.DeleteCollection(cfg.Name)
}

func deleteDeploymentConfig(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)
	return kconfig.Delete(cfg.Name)
}
