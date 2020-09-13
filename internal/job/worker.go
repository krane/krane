package job

import (
	"os"
	"time"

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
			job.startedAt = time.Now().Unix()
			job.state = InProgress

			err := job.Run(job.Args)
			if err != nil {
				logrus.Errorf("Error proceesing job %s", err.Error())
			}

			job.state = Complete
			job.completedAt = time.Now().Unix()
		case <-w.quit:
			logrus.Debug("Quitting worker")
			return
		}
	}
}
