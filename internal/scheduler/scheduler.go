package scheduler

import (
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/docker"
	job "github.com/biensupernice/krane/internal/jobs"
	"github.com/biensupernice/krane/internal/kranecfg"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type Scheduler struct {
	Store        store.Store
	dockerClient *client.Client
	Enqueuer     job.Enqueuer

	interval time.Duration
}

func New(store store.Store, dockerClient *client.Client, jobEnqueuer job.Enqueuer, interval_ms string) Scheduler {
	ms, _ := time.ParseDuration(interval_ms + "ms")
	return Scheduler{store, dockerClient, jobEnqueuer, ms}
}

func (sc *Scheduler) Run() {
	logrus.Infof("Starting Scheduler")

	for {
		sc.poll()
		<-time.After(sc.interval)
	}

	logrus.Infof("Exiting Scheduler")
	return
}

func (s *Scheduler) poll() {
	logrus.Debugf("Scheduler polling")

	// get docker containers
	containers, err := docker.GetKraneManagedContainers()
	if err != nil {
		logrus.Error("Scheduler error: %s", err)
		return
	}
	logrus.Debugf("Got %d containers", len(containers))

	// map deployments to containers
	deployments := mapDeploymentsToContainers(s.deployments(), containers)
	for name, obj := range deployments {
		if hasDesiredState(obj.config, obj.containers) {
			continue
		}

		// Serialize the desired state (config) to pass to the Job handler
		var args map[string]interface{}
		bytes, _ := json.Marshal(obj.config)
		json.Unmarshal(bytes, &args)

		// Persist job with status pending
		job := job.Job{
			ID:         utils.MakeIdentifier(),
			Namespace:  name,
			EnqueuedAt: time.Now().Unix(),
			Args:       args,

			Run: func(args job.Args) error {
				deployment := args["name"]
				logrus.Infof("Scheduler Handler: %s", deployment)
				return nil
			},
		}

		go s.Enqueuer.Enqueue(job)
	}

	logrus.Debugf("Next poll in %s", s.interval.String())
}

func hasDesiredState(kcfg kranecfg.KraneConfig, containers []types.Container) bool {
	return false
}

func (sc *Scheduler) deployments() []kranecfg.KraneConfig {
	bytes, err := sc.Store.GetAll(collection.Deployments)
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
