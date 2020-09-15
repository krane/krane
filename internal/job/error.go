package job

type Error struct {
	Message string `json:"message"`

}

func (job *Job) CaptureError(err error) {
	job.Status.Failures = append(job.Status.Failures, Error{err.Error()})
}
