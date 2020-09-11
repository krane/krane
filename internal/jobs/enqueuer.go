package job

import (
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/store"
)

const QueuedJobsCollections = "queued_jobs"

type Enqueuer struct {
	store store.Store

	jobQueue chan Job

	Handler GenericHandler
}

func NewEnqueuer(store store.Store, jobQueue chan Job) Enqueuer {
	return Enqueuer{store: store, jobQueue: jobQueue, Handler: nil}
}

func (e *Enqueuer) Enqueue(job Job) (Job, error) {
	logrus.Debugf("Queueing new job %s", job.ID)
	e.jobQueue <- job
	logrus.Debugf("Job %s Queued", job.ID)

	return job, nil
}
