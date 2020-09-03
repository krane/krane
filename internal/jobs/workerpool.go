package job

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/utils"
)

type WorkerPool struct {
	workerPoolID string
	concurrency  uint
	started      bool

	store *storage.Storage

	workers []*worker
}

// NewWorkerPool : create a concurrent pool of workers
func NewWorkerPool(concurrency uint, c *chan Job, store *storage.Storage) *WorkerPool {
	logrus.Debugf("Bootstrapping new worker pool with %d worker(s)", concurrency)
	wp := &WorkerPool{
		workerPoolID: utils.MakeIdentifier(),
		concurrency:  concurrency,
		store:        store,
	}

	for i := uint(0); i < wp.concurrency; i++ {
		logrus.Debugf("[wp %s] Appending new worker to the worker pool", wp.workerPoolID)
		w := newWorker(wp.workerPoolID)
		wp.workers = append(wp.workers, w)
	}

	return wp
}

// Start : all the workers part of the worker pool
func (wp *WorkerPool) Start() {
	if wp.started {
		return
	}

	wp.started = true

	for _, w := range wp.workers {
		logrus.Debugf("[wp %s][w %s] Starting new worker in separate thread", w.poolID, w.workerID)

		// go w.start()
		w.start()
	}
}

func (wp *WorkerPool) Stop() {
	if !wp.started {
		return
	}

	wp.started = false

	wg := sync.WaitGroup{}

	for _, w := range wp.workers {
		wg.Add(1)
		go func(w *worker) {
			// w.stop()
			wg.Done()
		}(w)
	}

	wg.Wait()

	logrus.Debugf("Stopping Worker Pool: %s", w.workerPoolID)
}
