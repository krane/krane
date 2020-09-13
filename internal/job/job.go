package job

import (
	"encoding/json"
)

type Job struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Args      Args   `json:"args"`

	state       State `json:"state"`
	enqueuedAt  int64 `json:"enqueued_at"`
	startedAt   int64 `json:"started_at"`
	completedAt int64 `json:"completed_at"`

	Run GenericHandler `json:"-"`
}

// Args : is a shortcut to easily specify arguments for job when enqueueing them.
type Args map[string]interface{}

// GenericHandler is a job handler without any custom context.
type GenericHandler func(Args) error

func (j *Job) serialize() ([]byte, error) {
	return json.Marshal(j)
}
