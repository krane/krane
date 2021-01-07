package job

import (
	"os"

	"github.com/pkg/errors"

	"github.com/krane/krane/internal/logger"
)

type worker struct {
	workerPool chan chan Job
	channel    chan Job
	quit       chan bool
}

// newWorker is a helper for creating new workers; a worker runs in its
// own routine waiting to process work from a job queue
func newWorker(workerPool chan chan Job, jobChannel chan Job) *worker {
	return &worker{workerPool, jobChannel, make(chan bool)}
}

// Start starts a worker
func (w *worker) start() {
	logger.Debugf("Worker starting with pid: %d", os.Getpid())
	go w.loop()
}

// stop stops a worker
func (w *worker) stop() {
	logger.Debug("Worker stopping")
	w.quit <- true
	return
}

// loop will infinitely block for jobs to come through from job queue
func (w *worker) loop() {
	logger.Debug("Worker loop started")
	for {
		select {
		case job := <-w.channel:
			job.start()

			for i := 0; i < int(job.RetryPolicy); i++ {
				job.Status.ExecutionCount++

				if job.Setup != nil {
					logger.Debugf("Setting up job %s", job.ID)
					if err := job.Setup(job.Args); err != nil {
						job.WithError(err)
						job.Status.FailureCount++
						continue
					}
				}

				if job.Run == nil {
					job.WithError(errors.New("job must have a Run implementation"))
					job.Status.FailureCount++
					return
				}

				if err := job.Run(job.Args); err != nil {
					job.WithError(err)
					job.Status.FailureCount++
					continue
				}

				if job.Finally != nil {
					logger.Debugf("Tearing down job %s", job.ID)
					if err := job.Finally(job.Args); err != nil {
						job.WithError(err)
						job.Status.FailureCount++
						continue
					}
				}

				logger.Debugf("Completed job %s", job.ID)
			}

			job.end()
		case <-w.quit:
			logger.Debug("Quitting worker")
			return
		}
	}
}
