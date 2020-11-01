package service

import (
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/store"
)

func StartDeployment(cfg config.Kconfig) error {
	j, err := makeDockerDeploymentJob(cfg, Up)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(j)
	return nil
}

func DeleteDeployment(cfg config.Kconfig) error {
	j, err := makeDockerDeploymentJob(cfg, Down)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(j)
	return nil
}

func enqueueDeploymentJob(deploymentJob job.Job) {
	db := store.Instance()
	queue := job.GetJobQueue()

	enq := job.NewEnqueuer(db, queue)
	_, err := enq.Enqueue(deploymentJob)
	if err != nil {
		logrus.Errorf("Error enqueuing deployment job %s", err.Error())
		return
	}
}
