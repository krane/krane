package scheduler

import (
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/kranecfg"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type Scheduler struct {
	store    store.Store
	docker   *docker.DockerClient
	enqueuer job.Enqueuer

	interval time.Duration
}

func New(store store.Store, dockerClient *docker.DockerClient, jobEnqueuer job.Enqueuer, interval_ms string) Scheduler {
	ms, _ := time.ParseDuration(interval_ms + "ms")
	return Scheduler{store, dockerClient, jobEnqueuer, ms}
}

func (sc *Scheduler) Run() {
	logrus.Debugf("Starting Scheduler")

	for {
		go sc.poll()
		<-time.After(sc.interval)
	}

	logrus.Debugf("Exiting Scheduler")
	return
}

func (s *Scheduler) poll() {
	logrus.Debugf("Scheduler polling")
	for _, deployment := range s.deployments() {
		containers, err := s.docker.FilterContainersByDeployment(deployment.Name)
		if err != nil {
			logrus.Errorf("Unhandled error when polling: %s", err)
			continue
		}

		if hasDesiredState(deployment, containers) {
			continue
		}

		// Serialize the deployment into a generic interface to pass as args to the Job handler
		var args map[string]interface{}
		bytes, _ := deployment.Serialize()
		json.Unmarshal(bytes, &args)

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

func hasDesiredState(kcfg kranecfg.KraneConfig, containers []types.Container) bool {
	return false
}

func (sc *Scheduler) deployments() []kranecfg.KraneConfig {
	bytes, err := sc.store.GetAll(collection.Deployments)
	if err != nil {
		logrus.Error("Scheduler error: %s", err)
		return make([]kranecfg.KraneConfig, 0)
	}

	deployments := make([]kranecfg.KraneConfig, 0)
	for _, b := range bytes {
		var d kranecfg.KraneConfig
		err := store.Deserialize(b, &d)
		if err != nil {
			logrus.Error("Unable to deserialize krane config", err.Error())
		}

		deployments = append(deployments, d)
	}
	return deployments
}
