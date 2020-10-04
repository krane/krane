package service

import (
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/store"
)

const retryPolicy = 3

type action string

func StartDeployment(cfg config.Config) error {
	j, err := makeDockerDeploymentJob(cfg, Up)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(j)
	return nil
}

func DeleteDeployment(cfg config.Config) error {
	j, err := makeDockerDeploymentJob(cfg, Down)
	if err != nil {
		return err
	}
	go enqueueDeploymentJob(j)
	return nil
}

func enqueueDeploymentJob(deploymentJob job.Job) {
	store := store.Instance()
	queue := job.GetJobQueue()

	e := job.NewEnqueuer(store, queue)
	_, err := e.Enqueue(deploymentJob)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
}
