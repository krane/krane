package bbq

import (
	"github.com/docker/distribution/uuid"
)

var (
	MaxWorker = 1 // os.Getenv("MAX_WORKERS")
	MaxQueue  = 1 // os.Getenv("MAX_QUEUE")
)

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	ID         string
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		ID:         uuid.Generate().String(),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {

	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Metadata.WorkerID = w.ID

				// Start processing job
				err := job.Process(job.Props)
				if err != nil {
					job.OnError(&job, err)
				}

				// Job has finished processing
				err = job.Done(&job)
				if err != nil {
					job.OnError(&job, err)
				}
			case <-w.quit:
				// we have received a signal to stop
				w.Stop()
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
