package service

import (
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
)

func applyDeployment(cfg config.DeploymentConfig, action DeploymentAction) error {
	deploymentJob, err := createDeploymentJob(cfg, action)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(deploymentJob)
	return nil
}

// enqueueDeploymentJob enqueues a deployment job for processing
func enqueueDeploymentJob(deploymentJob job.Job) {
	enqueuer := job.NewEnqueuer(job.Queue())
	queuedJob, err := enqueuer.Enqueue(deploymentJob)
	if err != nil {
		logger.Errorf("Error enqueuing deployment job %v", err)
		return
	}
	logger.Debugf("Queued job for %s", queuedJob.Namespace)
}

// StartDeployment starts a deployments container resources
func StartDeployment(cfg config.DeploymentConfig) error {
	return applyDeployment(cfg, CreateContainers)
}

// DeleteDeployment delete a deployment and its container resources
func DeleteDeployment(cfg config.DeploymentConfig) error {
	return applyDeployment(cfg, DeleteContainers)
}

// StopDeployment stops a deployments container resources
func StopDeployment(cfg config.DeploymentConfig) error {
	return applyDeployment(cfg, StopContainers)
}
