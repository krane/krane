package job

type Error struct {
	Execution uint   `json:"execution"`
	Message   string `json:"message"`
}

// WithError adds an error to a job
func (j *Job) WithError(err error) {
	j.Status.Failures = append(j.Status.Failures, Error{j.Status.ExecutionCount, err.Error()})
}
