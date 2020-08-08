package jobs

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/service/activity"
	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/pkg/bbq"
)

var (
	JobsCollectionName = "jobs"
)

func onJobError(job *bbq.Job, err error) {
	job.Success = false
	job.Error = err.Error()
}

func onJobDone(job *bbq.Job) error {
	if job.Error == "" {
		job.Success = true
	}

	// Capture the job activity
	jobActivity := &activity.Activity{Job: *job}
	activity.Capture(jobActivity)

	return removeJob(*job)
}

func storeJob(job bbq.Job) error {
	bytes, err := json.Marshal(job)
	if err != nil {
		logrus.Errorf("[%s] -> Error %s", job.ID, err.Error())
		return err
	}

	// Store job
	err = storage.Put(JobsCollectionName, job.ID, bytes)
	if err != nil {
		logrus.Errorf("[%s] -> Error %s", job.ID, err.Error())
		return err
	}
	return nil
}

func removeJob(job bbq.Job) error { return storage.Remove(JobsCollectionName, job.ID) }

func GetRunningJobs() ([]bbq.Job, error) {
	// Find jobs
	bytes, err := storage.GetAll(JobsCollectionName)
	if err != nil {
		return make([]bbq.Job, 0), err
	}

	if bytes == nil {
		return make([]bbq.Job, 0), nil
	}

	var jobs []bbq.Job
	for _, job := range bytes {
		var j bbq.Job
		err = json.Unmarshal(job, &j)
		if err != nil {
			return make([]bbq.Job, 0), err
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}
