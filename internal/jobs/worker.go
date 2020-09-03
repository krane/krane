package job

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/utils"
)

type Worker struct {
	ID         string
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

type worker struct {
	workerID string
	poolID   string
}

type WorkerPool struct {
	workerPoolID string
	concurrency  uint
	started      bool

	workers []*worker
}

// NewWorkerPool : create a concurrent pool of workers
func NewWorkerPool(concurrency uint) *WorkerPool {
	logrus.Debugf("Bootstrapping new worker pool with %d worker(s)", concurrency)
	wp := &WorkerPool{
		workerPoolID: utils.MakeIdentifier(),
		concurrency:  concurrency,
	}

	for i := uint(0); i < wp.concurrency; i++ {
		logrus.Debugf("[wp %s] Appending new worker to the worker pool", wp.workerPoolID)
		w := newWorker(wp.workerPoolID)
		wp.workers = append(wp.workers, w)
	}

	return wp
}

// newWorker : Helper for creating new workers
func newWorker(poolID string) *worker {
	workerID := utils.MakeIdentifier()

	w := &worker{
		poolID:   poolID,
		workerID: workerID,
	}

	return w
}

// func NewWorker(workerPool chan chan Job) Worker {
// 	return Worker{
// 		ID:         uuid.Generate().String(),
// 		WorkerPool: workerPool,
// 		JobChannel: make(chan Job),
// 		quit:       make(chan bool),
// 	}
// }

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

func (w *worker) start() {
	w.loop()
	// go w.loop()
	// go w.observer.start()
}

func (w *worker) loop() {
	logrus.Debugf("[wp %s][w %s] Starting worker loop", w.poolID, w.workerID)

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			logrus.Debugf("[wp %s][w %s] Processing new Jobs", w.poolID, w.workerID)
			time.Sleep(10 * time.Millisecond)
			timer.Reset(10 * time.Millisecond)
		}
	}
}

func (w *WorkerPool) Stop() {
	logrus.Debugf("Stopping Worker Pool: %s", w.workerPoolID)
}
