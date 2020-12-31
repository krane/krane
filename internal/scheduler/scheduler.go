package scheduler

import (
	"time"

	"github.com/pkg/errors"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
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
		containers, err := container.GetKraneContainersByDeployment(deployment.Name)
		if err != nil {
			logger.Error(errors.Wrap(err, "Unhandled error when polling"))
			continue
		}

		if hasDesiredState(deployment, containers) {
			continue
		}

		go s.enqueuer.Enqueue(job.Job{
			ID:         utils.ShortID(),
			Deployment: deployment.Name,
			Args:       map[string]interface{}{},
			Run: func(args interface{}) error {
				return nil
			},
		})
	}
	logger.Debugf("Next poll in %s", s.interval.String())
}

func hasDesiredState(kcfg config.DeploymentConfig, containers []container.KraneContainer) bool {
	// TODO: implementation not defined - always returning true to avoid doing anything
	return true
}

func (s *Scheduler) deployments() []config.DeploymentConfig {
	deployments, _ := config.GetAllDeploymentConfigurations()
	return deployments
}
