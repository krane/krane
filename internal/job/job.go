package job

import (
	"container/heap"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type Job struct {
	ID          string         `json:"id"`               // Unique job ID
	Namespace   string         `json:"namespace"`        // The namespace used for scoping jobs. This is the same namespace used when fetching secrets.
	Type        string         `json:"type"`             // The type of job
	Status      Status         `json:"status"`           // The status of the current job with details for execution counts etc..
	State       State          `json:"state"`            // Current state of a job (running | complete)
	StartTime   int64          `json:"start_time_epoch"` // Job start time - epoch in seconds since 1970
	EndTime     int64          `json:"end_time_epoch"`   // Job end time - epoch in seconds since 1970
	RetryPolicy uint           `json:"retry_policy"`     // Job retry policy
	Args        Args           `json:"-"`                // Arguments passed down to the Job Handler
	Run         GenericHandler `json:"-"`                // Executor function which receives the Args and returns an error if any
}

// Args : is a shortcut to easily specify arguments for job when enqueueing them.
type Args map[string]interface{}

// GenericHandler is a job handler without any custom context.
type GenericHandler func(Args) error

func (job *Job) serialize() ([]byte, error) { return json.Marshal(job) }

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
	job.State = Completed
	job.capture()
}

func (job *Job) capture() {
	// Unique constants to capture the jobs activity, format: {namespace}-JobsCollectionName
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

func CreateCollection(namespace string) error {
	collection := getNamespaceCollectionName(namespace)
	return store.Instance().CreateCollection(collection)
}

func DeleteCollection(namespace string) error {
	collection := getNamespaceCollectionName(namespace)
	return store.Instance().DeleteCollection(collection)
}

func (job *Job) validate() error {
	if job.ID == "" {
		return fmt.Errorf("job id required")
	}

	if job.Namespace == "" {
		return fmt.Errorf("job namespace required")
	}

	if job.Run == nil {
		return fmt.Errorf("unknown job handler")
	}

	maxRetryPolicy := utils.GetUIntEnv("JOB_MAX_RETRY_POLICY")
	if job.RetryPolicy > maxRetryPolicy {
		return fmt.Errorf("retry policy %d exceeds max retry policy %d", job.RetryPolicy, maxRetryPolicy)
	}

	// Every job should run under a namespace (the deployment scope).
	// If a new job being created is not bounded to a namespace an error will be thrown.
	ok, err := job.hasExistingNamespace()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("invalid job, %s", err.Error())
	}

	return nil
}

func (job *Job) hasExistingNamespace() (bool, error) {
	deployments, err := store.Instance().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return false, fmt.Errorf("invalid job, %s", err.Error())
	}

	found := false
	for _, deployment := range deployments {
		var d kconfig.Kconfig
		err := store.Deserialize(deployment, &d)
		if err != nil {
			return false, fmt.Errorf("invalid job, %s", err.Error())
		}

		if job.Namespace == d.Name {
			found = true
		}
	}

	if !found {
		return false, errors.New("invalid job, namespace not found")
	}

	return true, nil
}

func GetJobs(daysAgo uint) ([]Job, error) {
	// get all deployments
	deployments, err := kconfig.GetAllDeploymentConfigs()
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

	return recentActivity.mergeAndSort(), nil
}

func GetJobByID(namespace, id string, daysAgo uint) (Job, error) {
	jobs, err := GetJobsByNamespace(namespace, daysAgo)
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

func GetJobsByNamespace(namespace string, daysAgo uint) ([]Job, error) {
	// get start, end time range to get jobs for
	minDate, maxDate := calculateTimeRange(int(daysAgo))

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

// K dimensional Job array
type kJobs [][]Job

// merge : combines kjobs into a single job array sorted in DESCENDING order based on timestamp.
func (kjobs kJobs) mergeAndSort() []Job {
	var JobHeap jobHeap
	heap.Init(&JobHeap)

	// flatten kjobs into a single un-sorted array
	var flattened []Job
	for i := 0; i < len(kjobs); i++ {
		flattened = append(flattened, kjobs[i]...)
	}

	if flattened == nil {
		return make([]Job, 0)
	}

	// push all job's into heap
	for i := 0; i < len(flattened); i++ {
		heap.Push(&JobHeap, flattened[i])
	}

	// TODO: document the rest
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

	// Add the last item in the heap
	result = append(result, temp)

	return result
}

func getNamespaceCollectionName(namespace string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", namespace, constants.JobsCollectionName))
}

func calculateTimeRange(daysAgo int) (string, string) {
	start := time.Now().AddDate(0, 0, -daysAgo).Format(time.RFC3339)
	end := time.Now().Local().Format(time.RFC3339)

	return start, end
}
