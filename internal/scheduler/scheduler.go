package scheduler

import (
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
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
	logger.Debug("Starting Scheduler")

	for {
		go s.poll()
		<-time.After(s.interval)
	}

	logger.Debug("Exiting Scheduler")
}

func (s *Scheduler) poll() {
	logger.Debug("Scheduler polling")
	for _, deployment := range s.deployments() {
		containers, err := s.docker.FilterContainersByDeployment(deployment.Name)
		if err != nil {
			logger.Error(errors.Wrap(err, "Unhandled error when polling"))
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
				logger.Infof("Scheduler Handler: %s", args["name"])
				return nil
			},
		}

		go s.enqueuer.Enqueue(job)
	}
	logger.Debugf("Next poll in %s", s.interval.String())
}

func hasDesiredState(kcfg config.DeploymentConfig, containers []types.ContainerJSON) bool {
	// TODO: implementation not defined - always returning true to avoid doing anything
	return true
}

func (s *Scheduler) deployments() []config.DeploymentConfig {
	bytes, err := s.store.GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		logger.Errorf("Scheduler error: %s", err)
		return make([]config.DeploymentConfig, 0)
	}

	deployments := make([]config.DeploymentConfig, 0)
	for _, b := range bytes {
		var d config.DeploymentConfig
		err := store.Deserialize(b, &d)
		if err != nil {
			logger.Errorf("Unable to deserialize krane config", err)
		}

		deployments = append(deployments, d)
	}
	return deployments
}
