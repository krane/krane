package service

import (
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
)

// StartDeployment : starts a deployment
func StartDeployment(cfg config.DeploymentConfig) error {
	deploymentJob, err := createDeploymentJob(cfg, Up)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(deploymentJob)
	return nil
}

// DeleteDeployment delete a deployment
func DeleteDeployment(cfg config.DeploymentConfig) error {
	deploymentJob, err := createDeploymentJob(cfg, Down)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(deploymentJob)
	return nil
}

// enqueueDeploymentJob : enqueue a deployment job for processing
func enqueueDeploymentJob(deploymentJob job.Job) {
	enqueuer := job.NewEnqueuer(job.Queue())
	queuedJob, err := enqueuer.Enqueue(deploymentJob)
	if err != nil {
		logger.Errorf("Error enqueuing deployment job %v", err)
		return
	}
	logger.Debugf("Queued job for %s", queuedJob.Namespace)
}
