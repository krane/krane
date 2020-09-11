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
		case job := <-w.jobChannel:
			logrus.Infof("Got job for %s", job.Namespace)
			logrus.Debugf("Got job for %s", job.Namespace)
			job.Run(job.Args)
		case <-w.quit:
			logrus.Debug("Quitting worker")
			return
		}
	}
}
