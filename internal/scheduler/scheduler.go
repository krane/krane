package scheduler

import (
	"context"
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

// 1. TODO: Fetch current state
// 2. TODO: Fetch desired state
// 3. TODO: Queue required jobs
func (sc *Scheduler) Run() {
	logrus.Infof("Starting Scheduler")
	logrus.Debugf("Polling every %s", sc.interval.String())

	sc.Enqueuer.WithHandler("DEPLOY_JOB", handleRunDeployment)

	for {
		sc.poll(sc.interval)
		<-time.After(sc.interval)
	}

	// sc.Enqueuer.Enqueue("DEPLOY_JOB", job.Args{"id": status.MakeIdentifier()})
	// sc.Enqueuer.Enqueue("DELETE_JOB", job.Args{"id": status.MakeIdentifier()})
	return
}

func (sc *Scheduler) poll(interval_ms time.Duration) {
	logrus.Infof("Scheduler polling")

	// get deployments
	deployments := sc.deployments()
	logrus.Debugf("Got %d deployment", len(deployments))

	// get docker containers
	ctx := context.Background()
	kcontainers, err := docker.GetKraneManagedContainers(&ctx)
	if err != nil {
		logrus.Error("Scheduler error: %s", err)
		return
	}
	ctx.Done()

	if len(kcontainers) == 0 {
		logrus.Debugf("Found 0 krane managed containers. Waiting for next poll.")
	}

	logrus.Debugf("Got %d krane managed containers", len(kcontainers))

	// map deployments to containers
	mapping := mapDeploymentsToContainers(deployments, kcontainers)
	redeploymentJobs := make([]job.Args, 0)
	for _, state := range mapping {
		shouldRedeploy := !hasDesiredState(state.desiredState, kcontainers)

		// TODO: check if pending job already exists

		if shouldRedeploy {
			var args map[string]interface{}
			bytes, _ := json.Marshal(state)
			json.Unmarshal(bytes, &args)

			redeploymentJobs = append(redeploymentJobs, args)
		}
	}

	for _, jobArgs := range redeploymentJobs {
		sc.Enqueuer.Enqueue("DEPLOY_JOB", jobArgs)
	}

	logrus.Debugf("Next poll in %s", sc.interval.String())
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

func hasDesiredState(kcfg kranecfg.KraneConfig, containers []types.Container) bool {
	logrus.Infof("%v %v", kcfg, containers)
	return false
}

func handleRunDeployment(args job.Args) error {
	logrus.Infof("Running deployment job %v", args)
	return nil
}

func handleDeleteDeployment(args job.Args) error {
	logrus.Infof("Running deletion job %v", args)
	return nil
}
