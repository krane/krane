package job

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/utils"
)

type WorkerPool struct {
	workerPoolID string

	concurrency uint

	started bool

	store *storage.Storage

	workers    []*worker
	workerPool chan chan Job
	jobChannel chan Job
}

// NewWorkerPool : create a concurrent pool of workers to process Jobs from the jobQueue
func NewWorkerPool(concurrency uint, jobChannel chan Job, store *storage.Storage) WorkerPool {
	logrus.Debugf("Creating new worker pool with %d worker(s)", concurrency)
	wpID := utils.MakeIdentifier()
	wp := WorkerPool{
		workerPoolID: wpID,
		concurrency:  concurrency,
		store:        store,
		workerPool:   make(chan chan Job, concurrency),
		jobChannel:   jobChannel,
	}

	for i := uint(0); i < wp.concurrency; i++ {
		logrus.Infof("Appending new worker to worker pool %s", wp.workerPoolID)
		w := newWorker(wp.workerPool, wp.jobChannel)
		wp.workers = append(wp.workers, w)
	}

	logrus.Infof("%d worker(s) in the worker pool", len(wp.workers))

	return wp
}

// Start : all the workers part of the worker pool
func (wp *WorkerPool) Start() {
	logrus.Infof("Worker pool started on pid: %d", os.Getppid())
	if wp.started {
		return
	}

	wp.started = true

	var workersStarted int
	for _, w := range wp.workers {
		logrus.Info("Starting new go routine for worker")
		w.start()
		workersStarted++
	}

	logrus.Infof("Started %d worker(s)", workersStarted)

	return
}

func (wp *WorkerPool) Stop() {
	logrus.Infof("Stopping worker pool %s", wp.workerPoolID)

	if !wp.started {
		logrus.Info("Worker pool can't stop, it has not started")
		return
	}

	wp.started = false

	var workersStopped int

	var wg sync.WaitGroup
	for _, w := range wp.workers {
		logrus.Info("Adding worker to waitgroup")
		wg.Add(1)
		go func(w *worker) {
			logrus.Info("Attempting to stop Worker")
			w.stop()
			wg.Done()
			logrus.Info("Worker stopped, removing from waitgroup")
		}(w)
		workersStopped++
	}

	wg.Wait()
	logrus.Infof("%d worker(s) stopped", workersStopped)
}
