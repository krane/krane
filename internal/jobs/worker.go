package job

import (
	"os"

	"github.com/sirupsen/logrus"
)

type worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

// newWorker : Helper for creating new workers
func newWorker(workerPool chan chan Job, jobChannel chan Job) *worker {
	return &worker{workerPool, jobChannel, make(chan bool)}
}

func (w *worker) start() {
	logrus.Infof("Worker starting with pid: %d", os.Getpid())
	go w.loop()
}

func (w *worker) stop() {
	logrus.Info("Worker stopping")
	w.quit <- true
	return
}

func (w *worker) loop() {
	logrus.Infof("Worker loop started")
	for {
		select {
		case job := <-w.jobChannel:
			logrus.Infof("Got job %s", job.JobName)
			job.Run(job.Args)
		case <-w.quit:
			logrus.Info("Quitting worker")
			return
		}
	}
}
