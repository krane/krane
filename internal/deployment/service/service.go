package service

import (
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/job"
)

func StartDeployment(cfg kconfig.Kconfig) error {
	deploymentJob, err := makeDockerDeploymentJob(cfg, Up)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(deploymentJob)
	return nil
}

func DeleteDeployment(cfg kconfig.Kconfig) error {
	deploymentJob, err := makeDockerDeploymentJob(cfg, Down)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(deploymentJob)
	return nil
}

func enqueueDeploymentJob(deploymentJob job.Job) {
	queue := job.GetJobQueue()
	enqueuer := job.NewEnqueuer(queue)
	queuedJob, err := enqueuer.Enqueue(deploymentJob)
	if err != nil {
		logrus.Errorf("Error enqueuing deployment job for %s, %v", deploymentJob.Namespace, err)
		return
	}
	logrus.Debugf("Queued job for %s", queuedJob.Namespace)
}
