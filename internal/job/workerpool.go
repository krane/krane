package job

import (
	"os"
	"sync"

	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/store"
	"github.com/krane/krane/internal/utils"
)

type WorkerPool struct {
	workerPoolID string

	concurrency uint

	started bool

	store store.Store

	workers    []*worker
	workerPool chan chan Job
	jobChannel chan Job
}

// NewWorkerPool : create a concurrent pool of workers to process Jobs from the queue
func NewWorkerPool(concurrency uint, jobChannel chan Job, store store.Store) WorkerPool {
	logger.Debugf("Creating new worker pool with %d worker(s)", concurrency)
	wpID := utils.ShortID()
	wp := WorkerPool{
		workerPoolID: wpID,
		concurrency:  concurrency,
		store:        store,
		workerPool:   make(chan chan Job, concurrency),
		jobChannel:   jobChannel,
	}

	for i := uint(0); i < wp.concurrency; i++ {
		logger.Debugf("Appending new worker to worker pool %s", wp.workerPoolID)
		w := newWorker(wp.workerPool, wp.jobChannel)
		wp.workers = append(wp.workers, w)
	}

	logger.Debugf("%d worker(s) in the worker pool", len(wp.workers))

	return wp
}

// Start : all the workers part of the worker pool
func (wp *WorkerPool) Start() {
	logger.Debugf("Worker pool started on pid: %d", os.Getppid())
	if wp.started {
		return
	}

	wp.started = true

	var workersStarted int
	for _, w := range wp.workers {
		logger.Debug("Starting new worker")
		w.start()
		workersStarted++
	}

	logger.Debugf("Started %d worker(s)", workersStarted)

	return
}

func (wp *WorkerPool) Stop() {
	logger.Info("Disconnect signal received")

	if !wp.started {
		logger.Debug("Worker pool can't stop, it has not started")
		return
	}

	wp.started = false

	logger.Debugf("Stopping worker pool %s", wp.workerPoolID)

	stopped := 0
	var wg sync.WaitGroup
	for _, w := range wp.workers {
		logger.Debug("Adding worker to waitgroup")
		wg.Add(1)
		go func(w *worker) {
			logger.Debug("Attempting to stop Worker")
			w.stop()
			wg.Done()
			logger.Debug("Worker stopped, removing from waitgroup")
		}(w)
		stopped++
	}

	wg.Wait()
	logger.Debugf("%d out of %d worker(s) stopped", stopped, len(wp.workers))
}
