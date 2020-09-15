package job

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/kranecfg"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type Job struct {
	ID          string         `json:"id"`           // Unique job ID
	Namespace   string         `json:"namespace"`    // The namespace used for scoping jobs. This is the same namespace used when fetching secrets.
	Type        Type           `json:"type"`         // The type of job
	Status      Status         `json:"status"`       // The status of the current job with details for execution counts etc..
	State       State          `json:"state"`        // Current state of a job (running | complete)
	StartTime   int64          `json:"start_time"`   // Job start time - epoch in seconds since 1970
	EndTime     int64          `json:"end_time"`     // Job end time - epoch in seconds since 1970
	RetryPolicy uint           `json:"retry_policy"` // Job retry policy
	Args        Args           `json:"args"`         // Arguments passed down to the Job Handler
	Run         GenericHandler `json:"-"`            // Executor function which receives the Args and returns an error if any
}

// Args : is a shortcut to easily specify arguments for job when enqueueing them.
type Args map[string]interface{}

// GenericHandler is a job handler without any custom context.
type GenericHandler func(Args) error

func (job *Job) serialize() ([]byte, error) {
	return json.Marshal(job)
}

func (job *Job) start() {
	if job.State == Started {
		return
	}
	job.StartTime = time.Now().Unix()
	job.State = Started
}
func (job *Job) end() {
	if job.State != Started {
		return
	}
	job.EndTime = time.Now().Unix()
	job.State = Complete
	job.capture()
}

func (job *Job) capture() {
	// Unique collection to capture the jobs activity, format: {namespace}-Jobs
	collectionName := getNamespaceCollectionName(job.Namespace)
	bytes, _ := job.serialize()

	// timestamp is used as the key for the activity.
	// This leverages bolts time range scans which is an efficient way of performing lookups
	// for activity within a time range in an efficient manner.
	// The timestamp is RFC3339.
	timestamp := utils.UTCDateString()

	err := store.Instance().Put(collectionName, timestamp, bytes)
	if err != nil {
		logrus.Errorf("Unhandled error when inserting activity, %s", err)
		return
	}
}

func (job *Job) validate() error {
	if job.ID == "" {
		return fmt.Error("id required to create job")
	}

	if job.Namespace == "" {
		return fmt.Error("namespace required to create job")
	}

	if !isAllowedJobType(job.Type) {
		return fmt.Errorf("unknown job type %s", job.Type)
	}

	if job.Run == nil {
		return fmt.Errorf("unkown job handler")
	}

	maxRetryPolicy := utils.GetUIntEnv("JOB_MAX_RETRY_POLICY")
	if job.RetryPolicy > maxRetryPolicy {
		return fmt.Errorf("retry policy %d exceeds max retry policy %d", job.RetryPolicy, maxRetryPolicy)
	}

	return nil
}

func GetJobs(daysAgo uint) ([]Job, error) {
	// get all deployments
	deployments, err := kranecfg.GetAll()
	if err != nil {
		return make([]Job, 0), err
	}

	// K dim. arr containing un-merged deployment activities.
	var recentActivity kJobs = make([][]Job, 0)
	for _, deployment := range deployments {

		// get activity in time range
		deploymentActivity, err := GetJobsByNamespace(deployment.Name, daysAgo)
		if err != nil {
			return make([]Job, 0), err
		}

		recentActivity = append(recentActivity, deploymentActivity)
	}

	// Merging the activity orders all the deployment activity by time range
	return recentActivity.merge(), nil
}

func GetJobByID(namespace, id string) (Job, error) {
	jobs, err := GetJobsByNamespace(namespace, 365)
	if err != nil {
		return Job{}, fmt.Errorf("unable to fnd job with id %s", id)
	}

	for _, job := range jobs {
		if id == job.ID {
			return job, nil
		}
	}

	return Job{}, fmt.Errorf("unable to fnd job with id %s", id)
}

func GetJobsByNamespace(namespace string, daysAgo uint) ([]Job, error) {
	// get start, end time range
	minDate, maxDate := calculateTimeRange(daysAgo)

	// get activity in time range
	collectionName := getNamespaceCollectionName(namespace)
	bytes, err := store.Instance().GetInRange(collectionName, minDate, maxDate)
	if err != nil {
		return make([]Job, 0), err
	}

	recentActivity := make([]Job, 0)
	for _, activityBytes := range bytes {
		var j Job
		err := store.Deserialize(activityBytes, &j)
		if err != nil {
			return recentActivity, err
		}

		recentActivity = append(recentActivity, j)
	}
	return recentActivity, nil
}

// K dimension Job array with a method merge that combines into a single array and returns timestamp sorted jobs.
type kJobs [][]Job

func (kjobs kJobs) merge() []Job {
	jobs := make([]Job, 0)

	for _, k := range kjobs {
		jobs = append(jobs, sort(jobs, k)...)
	}

	return jobs
}

func sort(arr1 []Job, arr2 []Job) []Job {
	sorted := make([]Job, 0)

	i := 0
	j := 0

	for i < len(arr1) || j < len(arr2) {
		// Compare using start timestamp
		if arr1[i].StartTime < arr2[j].StartTime {
			sorted = append(sorted, arr1[i])
			i++
		} else {
			sorted = append(sorted, arr2[i])
			j++
		}
	}

	if i < len(arr1) {
		sorted = append(sorted, arr2...)
	}

	if j < len(arr2) {
		sorted = append(sorted, arr1...)
	}

	return sorted
}

func getNamespaceCollectionName(namespace string) string {
	return fmt.Sprintf("%s-%s", namespace, collection.Jobs)
}

func calculateTimeRange(daysAgo uint) (string, string) {
	start := time.Now().Add(time.Duration(24*daysAgo) * time.Hour).Format(time.RFC3339)
	end := time.Now().Local().Format(time.RFC3339)

	return start, end
}
