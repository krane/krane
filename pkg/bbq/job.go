package bbq

import (
	"github.com/docker/distribution/uuid"

	"github.com/biensupernice/krane/internal/utils"
)

type Job struct {
	ID        string            `json:"id"`
	CreatedAt string            `json:"created_at"`
	Body      interface{}       `json:"body"`
	Props     map[string]string `json:"props,omitempty"` // Config props passed down when processing the job
	JobType   string            `json:"job_type"`
	Metadata  Metadata          `json:"metadata"`
	Success   bool              `json:"success"`
	Error     string            `json:"error,omitempty"`

	Process func(props map[string]string) error `json:"-"`
	Done    func(job *Job) error                `json:"-"`
	OnError func(job *Job, err error)           `json:"-"`
}

type Metadata struct {
	WorkerID string `json:"worker_id"`
}

var JobQueue chan Job

func InitJobQueue() {
	// Create buffered job queue
	JobQueue = make(chan Job, MaxQueue)

	dispatcher := NewDispatcher(MaxWorker)
	dispatcher.Run()
}

func Queue(job Job) {
	job.ID = uuid.Generate().String()
	job.CreatedAt = utils.UTCDateString()

	JobQueue <- job
}
