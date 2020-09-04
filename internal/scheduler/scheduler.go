package scheduler

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"

	job "github.com/biensupernice/krane/internal/jobs"
	"github.com/biensupernice/krane/internal/storage"
)

type Scheduler struct {
	Store        storage.Storage
	dockerClient *client.Client
	Enqueuer     job.Enqueuer
}

func New(store storage.Storage, dockerClient *client.Client, jobEnqueuer job.Enqueuer) Scheduler {
	return Scheduler{store, dockerClient, jobEnqueuer}
}

// 1. TODO: Fetch current state
// 2. TODO: Fetch desired state
// 3. TODO: Queue required jobs
func (sc *Scheduler) Run() {
	logrus.Infof("Starting Scheduler")

	// sc.Enqueuer.WithHandler("DEPLOY_JOB", handleRunDeployment)
	// sc.Enqueuer.WithHandler("DELETE_JOB", handleDeleteDeployment)

	// sc.Enqueuer.Enqueue("DEPLOY_JOB", job.Args{"id": utils.MakeIdentifier()})
	// sc.Enqueuer.Enqueue("DELETE_JOB", job.Args{"id": utils.MakeIdentifier()})
	return
}

// func handleRunDeployment(args job.Args) error {
// 	logrus.Infof("Running deployment job %v", args)
// 	return nil
// }
//
// func handleDeleteDeployment(args job.Args) error {
// 	logrus.Infof("Running deletion job %v", args)
// 	return nil
// }
