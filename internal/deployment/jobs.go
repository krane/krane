package deployment

import (
	"fmt"

	"github.com/krane/krane/internal/job"
	"github.com/krane/krane/internal/store"
	"github.com/krane/krane/internal/utils"
)

// CreateCollection create the job collection for a deployment
func CreateJobsCollection(deployment string) error {
	collection := job.GetJobsCollectionName(deployment)
	return store.Client().CreateCollection(collection)
}

// DeleteCollection deletes the job collection for a deployment
func DeleteJobsCollection(deployment string) error {
	collection := job.GetJobsCollectionName(deployment)
	return store.Client().DeleteCollection(collection)
}

// GetJobs returns all deployment jobs
func GetJobs(daysAgo uint) ([]job.Job, error) {
	deployments, err := GetAllDeploymentConfigs()
	if err != nil {
		return make([]job.Job, 0), err
	}

	// N dim. arr containing un-merged deployment jobs
	recentActivity := make(job.NJobs, 0)
	for _, deployment := range deployments {
		activity, err := GetJobsByDeployment(deployment.Name, daysAgo)
		if err != nil {
			return make([]job.Job, 0), err
		}

		recentActivity = append(recentActivity, activity)
	}

	return recentActivity.MergeAndSort(), nil
}

// GetJobByID returns a job by id
func GetJobByID(deployment, id string, daysAgo uint) (job.Job, error) {
	jobs, err := GetJobsByDeployment(deployment, daysAgo)
	if err != nil {
		return job.Job{}, fmt.Errorf("unable to find a job with id %s", id)
	}

	for _, j := range jobs {
		if id == j.ID {
			return j, nil
		}
	}

	return job.Job{}, fmt.Errorf("unable to fnd job with id %s", id)
}

// GetJobs returns all jobs for a deployment within a time range
func GetJobsByDeployment(deployment string, daysAgo uint) ([]job.Job, error) {
	// get start & end dates for the range of jobs to look for
	minDate, maxDate := utils.CalculateTimeRange(int(daysAgo))

	// get activity in time range
	collection := job.GetJobsCollectionName(deployment)
	bytes, err := store.Client().GetInRange(collection, minDate, maxDate)
	if err != nil {
		return make([]job.Job, 0), err
	}

	recentActivity := make([]job.Job, 0)
	for _, activityBytes := range bytes {
		var j job.Job
		if err := store.Deserialize(activityBytes, &j); err != nil {
			return recentActivity, err
		}
		recentActivity = append(recentActivity, j)
	}
	return recentActivity, nil
}
