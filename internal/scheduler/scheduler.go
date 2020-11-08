package scheduler

import (
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type Scheduler struct {
	store    store.Store
	docker   *docker.Client
	enqueuer job.Enqueuer

	interval time.Duration
}

func New(store store.Store, dockerClient *docker.Client, jobEnqueuer job.Enqueuer, interval_ms string) Scheduler {
	ms, _ := time.ParseDuration(interval_ms + "ms")
	return Scheduler{store, dockerClient, jobEnqueuer, ms}
}

func (s *Scheduler) Run() {
	logrus.Debug("Starting Scheduler")

	for {
		go s.poll()
		<-time.After(s.interval)
	}

	logrus.Debug("Exiting Scheduler")
}

func (s *Scheduler) poll() {
	logrus.Debug("Scheduler polling")
	for _, deployment := range s.deployments() {
		containers, err := s.docker.FilterContainersByDeployment(deployment.Name)
		if err != nil {
			logrus.Error(errors.Wrap(err, "Unhandled error when polling"))
			continue
		}

		if hasDesiredState(deployment, containers) {
			continue
		}

		// Serialize the deployment into a generic interface to pass as args to the Job handler
		var args map[string]interface{}
		bytes, _ := deployment.Serialize()
		_ = json.Unmarshal(bytes, &args)

		job := job.Job{
			ID:        utils.MakeIdentifier(),
			Namespace: deployment.Name,
			Args:      args,
			Run: func(args job.Args) error {
				logrus.Infof("Scheduler Handler: %s", args["name"])
				return nil
			},
		}

		go s.enqueuer.Enqueue(job)
	}
	logrus.Debugf("Next poll in %s", s.interval.String())
}

func hasDesiredState(kcfg config.Kconfig, containers []types.ContainerJSON) bool {
	// TODO: implmentation not defined always returning true to avoid doing anything
	return true
}

func (s *Scheduler) deployments() []config.Kconfig {
	bytes, err := s.store.GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		logrus.Errorf("Scheduler error: %s", err)
		return make([]config.Kconfig, 0)
	}

	deployments := make([]config.Kconfig, 0)
	for _, b := range bytes {
		var d config.Kconfig
		err := store.Deserialize(b, &d)
		if err != nil {
			logrus.Error("Unable to deserialize krane config", err.Error())
		}

		deployments = append(deployments, d)
	}
	return deployments
}
