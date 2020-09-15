package job

import (
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/store"
)

type Enqueuer struct {
	store   store.Store
	queue   chan Job
	Handler GenericHandler
}

func NewEnqueuer(store store.Store, jobQueue chan Job) Enqueuer {
	return Enqueuer{store: store, queue: jobQueue, Handler: nil}
}

func (e *Enqueuer) Enqueue(job Job) (Job, error) {
	err := job.validate()
	if err != nil {
		return Job{}, err
	}

	logrus.Debugf("Queueing new job %s", job.ID)
	e.queue <- job
	logrus.Debugf("Job %s Queued", job.ID)
	return job, nil
}
