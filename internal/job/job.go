package job

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type Job struct {
	ID          string         `json:"id"`               // Unique job ID
	Deployment  string         `json:"deployment"`       // Deployment used for scoping jobs.
	Type        string         `json:"type"`             // The type of job
	Status      Status         `json:"response"`         // The response of the current job with details for execution counts etc..
	State       State          `json:"state"`            // Current state of a job (running | complete)
	StartTime   int64          `json:"start_time_epoch"` // Job Start time - epoch in seconds since 1970
	EndTime     int64          `json:"end_time_epoch"`   // Job end time - epoch in seconds since 1970
	RetryPolicy uint           `json:"retry_policy"`     // Job retry policy
	Args        interface{}    `json:"-"`                // Arguments passed down to job handlers
	Setup       GenericHandler `json:"-"`                // Setup is the initial execution fn for a job typically to setup arguments
	Run         GenericHandler `json:"-"`                // Run is the main executor fn for a job
	Finally     GenericHandler `json:"-"`                // Final fn is the final execution fn for a job
}

// GenericHandler is a generic job handler that takes in job arguments
type GenericHandler func(args interface{}) error

// Serialize a job into bytes
func (j *Job) Serialize() ([]byte, error) { return json.Marshal(j) }

// Start : Start a job
func (j *Job) start() {
	if j.State == Started {
		return
	}
	j.StartTime = time.Now().Unix()
	j.State = Started
}

func (j *Job) end() {
	if j.State != Started {
		return
	}
	j.EndTime = time.Now().Unix()
	j.State = Completed
	j.save()
}

// save : store the job
func (j *Job) save() {
	collection := getDeploymentJobsCollectionName(j.Deployment)
	bytes, _ := j.Serialize()

	// timestamp(RFC3339) is used as the key for the activity.
	// This leverages bolts time range scans which is an efficient way of performing lookups
	// for activity within a time range in an efficient manner.
	timestamp := utils.UTCDateString()

	err := store.Client().Put(collection, timestamp, bytes)
	if err != nil {
		logger.Errorf("Unhandled error when inserting j, %s", err)
		return
	}
}

// CreateCollection : create jobs collection for a deployment
func CreateCollection(namespace string) error {
	collection := getDeploymentJobsCollectionName(namespace)
	return store.Client().CreateCollection(collection)
}

// DeleteCollection : delete jobs collection for a deployment
func DeleteCollection(namespace string) error {
	collection := getDeploymentJobsCollectionName(namespace)
	return store.Client().DeleteCollection(collection)
}

// validate : validate a jobs configuration
func (j *Job) validate() error {
	if j.ID == "" {
		return fmt.Errorf("job id required")
	}

	if j.Deployment == "" {
		return fmt.Errorf("job deployment required")
	}

	if j.Run == nil {
		return fmt.Errorf("unknown job handler")
	}

	maxRetryPolicy := utils.UIntEnv(constants.EnvJobMaxRetryPolicy)
	if j.RetryPolicy > maxRetryPolicy {
		return fmt.Errorf("retry policy %d exceeds job max retry policy %d", j.RetryPolicy, maxRetryPolicy)
	}

	// every job should run under a deployment.
	// if a new job being created is not bounded to a deployment an error will be thrown.
	ok, err := j.hasExistingDeployment()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("invalid job, %v", err)
	}

	return nil
}

func (j *Job) hasExistingDeployment() (bool, error) {
	deployments, err := store.Client().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return false, fmt.Errorf("invalid job, %v", err)
	}

	for _, deployment := range deployments {
		var d config.DeploymentConfig
		if err := store.Deserialize(deployment, &d); err != nil {
			return false, fmt.Errorf("invalid job, %v", err)
		}
		if j.Deployment == d.Name {
			return true, nil
		}
	}

	return false, fmt.Errorf("invalid job, deployment %s not found", j.Deployment)
}

// GetJobs : get all jobs
func GetJobs(daysAgo uint) ([]Job, error) {
	// get all deployments
	deployments, err := config.GetAllDeploymentConfigurations()
	if err != nil {
		return make([]Job, 0), err
	}

	// K dim. arr containing un-merged deployment activities.
	var recentActivity nJobs = make([][]Job, 0)
	for _, deployment := range deployments {

		// get activity in time range
		deploymentActivity, err := GetJobsByDeployment(deployment.Name, daysAgo)
		if err != nil {
			return make([]Job, 0), err
		}

		recentActivity = append(recentActivity, deploymentActivity)
	}

	return recentActivity.mergeAndSort(), nil
}

// GetJobByID : get a job by id
func GetJobByID(namespace, id string, daysAgo uint) (Job, error) {
	jobs, err := GetJobsByDeployment(namespace, daysAgo)
	if err != nil {
		return Job{}, fmt.Errorf("unable to find a job with id %s", id)
	}

	for _, job := range jobs {
		if id == job.ID {
			return job, nil
		}
	}

	return Job{}, fmt.Errorf("unable to fnd job with id %s", id)
}

// GetJobs : get all jobs for a deployment
func GetJobsByDeployment(deployment string, daysAgo uint) ([]Job, error) {
	// get Start, end time range to get jobs for
	minDate, maxDate := calculateTimeRange(int(daysAgo))

	// get activity in time range
	collection := getDeploymentJobsCollectionName(deployment)
	bytes, err := store.Client().GetInRange(collection, minDate, maxDate)
	if err != nil {
		return make([]Job, 0), err
	}

	recentActivity := make([]Job, 0)
	for _, activityBytes := range bytes {
		var j Job
		if err := store.Deserialize(activityBytes, &j); err != nil {
			return recentActivity, err
		}
		recentActivity = append(recentActivity, j)
	}
	return recentActivity, nil
}

// N dimensional Job array
type nJobs [][]Job

// merge : combines njobs into a single job array sorted in DESCENDING order based on timestamp.
func (njobs nJobs) mergeAndSort() []Job {
	var JobHeap jobHeap
	heap.Init(&JobHeap)

	// flatten njobs into a single un-sorted array
	var flattened []Job
	for i := 0; i < len(njobs); i++ {
		flattened = append(flattened, njobs[i]...)
	}

	if flattened == nil {
		return make([]Job, 0)
	}

	// push all jobs into the heap
	for i := 0; i < len(flattened); i++ {
		heap.Push(&JobHeap, flattened[i])
	}

	overlap := func(a, b Job) bool {
		if a.StartTime > b.EndTime {
			return false
		}
		if b.StartTime > a.EndTime {
			return false
		}
		return true
	}

	temp := heap.Pop(&JobHeap).(Job)
	var result []Job

	for JobHeap.Len() > 0 {
		a := temp
		b := heap.Pop(&JobHeap).(Job)

		if overlap(a, b) {
			if a.EndTime < b.EndTime {
				temp = b
			}
		} else {
			result = append(result, a)
			temp = b
		}
	}

	// add the last item in the heap
	result = append(result, temp)

	return result
}

func getDeploymentJobsCollectionName(deployment string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", deployment, constants.JobsCollectionName))
}

func calculateTimeRange(daysAgo int) (string, string) {
	start := time.Now().AddDate(0, 0, -daysAgo).Format(time.RFC3339)
	end := time.Now().Local().Format(time.RFC3339)
	return start, end
}
