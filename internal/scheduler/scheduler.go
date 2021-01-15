package scheduler

import (
	"time"

	"github.com/pkg/errors"

	"github.com/krane/krane/internal/deployment"
	"github.com/krane/krane/internal/docker"
	"github.com/krane/krane/internal/job"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/store"
)

type Scheduler struct {
	store    store.Store
	docker   *docker.Client
	enqueuer job.Enqueuer
	interval time.Duration
}

// New returns a new scheduler used to poll and create deployment resources
func New(store store.Store, dockerClient *docker.Client, jobEnqueuer job.Enqueuer, interval_ms string) Scheduler {
	interval, _ := time.ParseDuration(interval_ms + "ms")
	return Scheduler{store, dockerClient, jobEnqueuer, interval}
}

// Run starts the scheduler polling on an interval
func (s *Scheduler) Run() {
	logger.Debug("Starting Scheduler")

	for {
		go s.poll()
		<-time.After(s.interval)
	}
}

// poll will on an interval get deployments and queue jobs if they are not
// in a desired state. For example, if a deployment has a scale of 3 but only
// 1 container is running, the scheduler schedules a new job to update the deployment state.
func (s *Scheduler) poll() {
	logger.Debug("Scheduler polling")

	deployments, err := deployment.GetAllDeployments()
	if err != nil {
		logger.Error(errors.Wrap(err, "Unhandled error when polling"))
	}

	for _, d := range deployments {
		if hasDesiredState(d) {
			continue
		}
	}

	logger.Debugf("Next poll in %s", s.interval.String())
}

// hasDesiredState checks that deployments are in parity with their configurations
func hasDesiredState(d deployment.Deployment) bool {
	config := d.Config
	containers := d.Containers

	if config.Scale != len(containers) {
		return false
	}

	if err := deployment.RetriableContainersHealthCheck(containers, 3); err != nil {
		return false
	}

	return true
}
