package service

import (
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
)

// StartDeployment:
func StartDeployment(cfg kconfig.Kconfig) error {
	deploymentJob, err := makeDockerDeploymentJob(cfg, Up)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(deploymentJob)
	return nil
}

// DeleteDeployment:
func DeleteDeployment(cfg kconfig.Kconfig) error {
	deploymentJob, err := makeDockerDeploymentJob(cfg, Down)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(deploymentJob)
	return nil
}

// enqueueDeploymentJob:
func enqueueDeploymentJob(deploymentJob job.Job) {
	queue := job.Queue()
	enqueuer := job.NewEnqueuer(queue)
	queuedJob, err := enqueuer.Enqueue(deploymentJob)
	if err != nil {
		logger.Errorf("Error enqueuing deployment job %v", err)
		return
	}
	logger.Debugf("Queued job for %s", queuedJob.Namespace)
}
