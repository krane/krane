package job

import (
	"encoding/json"
)

type Job struct {
	ID         string `json:"id"`
	JobName    string `json:"name"`
	EnqueuedAt int64  `json:"enqueued_at"`
	Args       Args   `json:"args"`

	Run GenericHandler `json:"-"`
}

// Args : is a shortcut to easily specify arguments for jobs when enqueueing them.
type Args map[string]interface{}

// GenericHandler is a job handler without any custom context.
type GenericHandler func(Args) error

func NewJobQueue(queueSize uint) chan Job {
	// logrus.Infof("Creating job queue of size %d", queueSize)
	return make(chan Job, queueSize)
}

func (j *Job) serialize() ([]byte, error) {
	return json.Marshal(j)
}
