package job

import (
	"encoding/json"
	"sync"

	"github.com/sirupsen/logrus"
)

type Job struct {
	ID         string `json:"id"`
	Namespace  string `json:"namespace"`
	EnqueuedAt int64  `json:"enqueued_at"`
	Args       Args   `json:"args"`

	Run GenericHandler `json:"-"`
}

// Args : is a shortcut to easily specify arguments for jobs when enqueueing them.
type Args map[string]interface{}

// GenericHandler is a job handler without any custom context.
type GenericHandler func(Args) error

var once sync.Once
var instance chan Job

func NewJobQueue(queueSize uint) chan Job {
	logrus.Debugf("Creating job queue of size %d", queueSize)
	once.Do(func() { instance = make(chan Job, queueSize) })
	return instance
}

func GetJobQueue() chan Job { return instance }

func (j *Job) serialize() ([]byte, error) {
	return json.Marshal(j)
}
