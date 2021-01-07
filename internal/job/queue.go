package job

import (
	"sync"

	"github.com/krane/krane/internal/logger"
)

var once sync.Once
var queue chan Job

// Queue : get the job queue
func Queue() chan Job { return queue }

// NewBufferedQueue : create a buffered channel for queuing jobs
func NewBufferedQueue(queueSize uint) chan Job {
	logger.Debugf("Creating job queue of size %d", queueSize)
	once.Do(func() { queue = make(chan Job, queueSize) })
	return queue
}
