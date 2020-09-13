package job

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var once sync.Once
var instance chan Job

func GetJobQueue() chan Job { return instance }

func NewJobQueue(queueSize uint) chan Job {
	logrus.Debugf("Creating job queue of size %d", queueSize)
	once.Do(func() { instance = make(chan Job, queueSize) })
	return instance
}
