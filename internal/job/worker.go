package job

import (
	"os"

	"github.com/biensupernice/krane/internal/logger"
)

type worker struct {
	workerPool chan chan Job
	channel    chan Job
	quit       chan bool
}

// newWorker : helper for creating new workers; a worker runs in its own go routine
// waiting for processing jobs from a the queue
func newWorker(workerPool chan chan Job, jobChannel chan Job) *worker {
	return &worker{workerPool, jobChannel, make(chan bool)}
}

// start : start a worker
func (w *worker) start() {
	logger.Debugf("Worker starting with pid: %d", os.Getpid())
	go w.loop()
}

// stop : stop a worker
func (w *worker) stop() {
	logger.Debug("Worker stopping")
	w.quit <- true
	return
}

// loop : a worker will wait in a blocking manner for jobs to come through the job queue.
func (w *worker) loop() {
	logger.Debug("Worker loop started")
	for {
		select {
		case job := <-w.channel:
			job.start()

			for i := 0; i < int(job.RetryPolicy); i++ {
				job.Status.ExecutionCount++
				err := job.Run(job.Args)
				if err == nil {
					logger.Debugf("Completed job %s for %s", job.ID, job.Namespace)
					break
				}
				logger.Errorf("Error processing job %v", err)
				job.WithError(err)
				job.Status.FailureCount++
			}
			job.end()
		case <-w.quit:
			logger.Debug("Quitting worker")
			return
		}
	}
}
