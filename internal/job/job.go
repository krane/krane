package job

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

type Job struct {
	ID          string         `json:"id"`               // Unique job ID
	Deployment  string         `json:"deployment"`       // Deployment used for scoping jobs.
	Type        string         `json:"type"`             // The type of job; Generic to allow for open types of jobs to be executed
	Status      Status         `json:"status"`           // The response of the current job with details for execution counts etc..
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
	collection := GetJobsCollectionName(j.Deployment)
	bytes, _ := j.Serialize()

	// timestamp(RFC3339) is used as the key for the activity.
	// This leverages bolts time range scans which is an efficient way of performing lookups
	// for activity within a time range in an efficient manner.
	timestamp := utils.UTCDateString()

	err := store.Client().Put(collection, timestamp, bytes)
	if err != nil {
		logger.Errorf("Unhandled error when inserting job into the db, %s", err)
		return
	}
}

// validate returns an error if a Job does not have a valid configuration
func (j *Job) validate() error {
	if j.ID == "" {
		return fmt.Errorf("job id required")
	}

	if j.Deployment == "" {
		return fmt.Errorf("deployment required to execute job")
	}

	if j.Run == nil {
		return fmt.Errorf("run must be implemented for a job")
	}

	maxRetryPolicy := utils.UIntEnv(constants.EnvJobMaxRetryPolicy)
	if j.RetryPolicy > maxRetryPolicy {
		return fmt.Errorf("retry policy %d exceeds job max retry policy %d", j.RetryPolicy, maxRetryPolicy)
	}

	return nil
}

// N dimensional Job array
type NJobs [][]Job

// merge : combines njobs into a single job array sorted in DESCENDING order based on timestamp.
func (njobs NJobs) MergeAndSort() []Job {
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

func GetJobsCollectionName(deployment string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", deployment, constants.JobsCollectionName))
}
