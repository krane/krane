package job

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/utils"
)

const QueuedJobsCollections = "queued_jobs"

type Handlers map[string]GenericHandler

type Enqueuer struct {
	store *storage.Storage

	jobQueue chan Job

	Handlers Handlers
}

func NewEnqueuer(store *storage.Storage, jobQueue chan Job) Enqueuer {
	return Enqueuer{store: store, jobQueue: jobQueue, Handlers: make(Handlers, 0)}
}

func (e *Enqueuer) Enqueue(jobName string, args Args) (Job, error) {
	logrus.Infof("Enqueueing new job")
	jobHandler := e.Handlers[jobName]
	if jobHandler == nil {
		logrus.Info("Unable to queue job, unknown handler")
		return Job{}, errors.New("unable to queue job, unknown handler")
	}

	job := Job{
		ID:         utils.MakeIdentifier(),
		EnqueuedAt: time.Now().Unix(),
		JobName:    jobName,
		Args:       args,

		Run: jobHandler,
	}

	bytes, err := job.serialize()
	if err != nil {
		return Job{}, err
	}

	logrus.Infof("Adding job %s to the store with status pending", job.ID)
	store := *e.store
	err = store.Put(QueuedJobsCollections, job.ID, bytes)
	if err != nil {
		return Job{}, err
	}
	logrus.Infof("Job %s added to the store with status pending", job.ID)

	logrus.Infof("Queueing new job %s", job.ID)
	e.jobQueue <- job
	logrus.Infof("Job %s Queued", job.ID)

	return job, nil
}

func (e *Enqueuer) WithHandler(jobName string, handler GenericHandler) {
	if jobName == "" {
		logrus.Info("Unable to register job handler, missing jobName")
		return
	}

	if handler == nil {
		logrus.Info("Unable to register job handler, missing job handler")
		return
	}

	logrus.Infof("Registering new job handler %s", jobName)
	e.Handlers[jobName] = handler
	logrus.Infof("Successfully registered new job handler %s", jobName)
}
