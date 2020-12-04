package job

import "github.com/biensupernice/krane/internal/logger"

type Enqueuer struct {
	queue   chan Job
	Handler GenericHandler
}

func NewEnqueuer(jobQueue chan Job) Enqueuer {
	return Enqueuer{queue: jobQueue, Handler: nil}
}

func (e *Enqueuer) Enqueue(job Job) (Job, error) {
	err := job.validate()
	if err != nil {
		return Job{}, err
	}

	logger.Debugf("Queueing new job %s", job.ID)
	e.queue <- job // Blocks here until space opens up in the queue
	logger.Debugf("Job %s Queued", job.ID)
	return job, nil
}
