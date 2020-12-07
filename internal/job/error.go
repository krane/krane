package job

type Error struct {
	Execution uint   `json:"execution"`
	Message   string `json:"message"`
}

// WithError : add an error to a job
func (job *Job) WithError(err error) {
	job.Status.Failures = append(job.Status.Failures, Error{job.Status.ExecutionCount, err.Error()})
}
