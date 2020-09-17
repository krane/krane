package job

import (
	"os"

	"github.com/sirupsen/logrus"
)

type worker struct {
	workerPool chan chan Job
	channel    chan Job
	quit       chan bool
}

// newWorker : Helper for creating new workers
func newWorker(workerPool chan chan Job, jobChannel chan Job) *worker {
	return &worker{workerPool, jobChannel, make(chan bool)}
}

func (w *worker) start() {
	logrus.Debugf("Worker starting with pid: %d", os.Getpid())
	go w.loop()
}

func (w *worker) stop() {
	logrus.Debug("Worker stopping")
	w.quit <- true
	return
}

func (w *worker) loop() {
	logrus.Debug("Worker loop started")
	for {
		select {
		case job := <-w.channel:
			job.start()

			for i := 0; i < int(job.RetryPolicy); i++ {
				job.Status.ExecutionCount++
				err := job.Run(job.Args)
				if err == nil {
					break
				}
				logrus.Errorf("Error procesing job %s", err.Error())
				job.WithError(err)
				job.Status.FailureCount++
			}

			job.end()
		case <-w.quit:
			logrus.Debug("Quitting worker")
			return
		}
	}
}