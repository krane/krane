package job

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/utils"
)

const PendingJobsCollections = "PJOBS"
const QueuedJobsCollections = "QJOBS"
const CompletedJobsCollection = "CJOBS"

type enqueuer struct {
	store *storage.Storage

	c *chan Job
}

func NewEnqueuer(store *storage.Storage, c *chan Job) *enqueuer {
	return &enqueuer{store, c}
}

func (e *enqueuer) Enqueue(jobName string, args map[string]interface{}) (*Job, error) {
	job := &Job{
		ID:         utils.MakeIdentifier(),
		EnqueuedAt: time.Now().Unix(),
		Name:       jobName,
		Args:       args,
	}

	bytes, err := job.serialize()
	if err != nil {
		return nil, err
	}

	logrus.Debugf("[%s %s] Adding job with status pending", job.Name, job.ID)
	store := *e.store
	err = store.Put(PendingJobsCollections, job.ID, bytes)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[%s %s] Job added with status pending", job.Name, job.ID)

	logrus.Debugf("[%s %s] Sending to Job channel", job.Name, job.ID)
	c := *e.c
	c <- *job
	logrus.Debugf("[%s %s] Job sent to channel", job.Name, job.ID)

	return job, nil
}
