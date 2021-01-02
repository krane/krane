package job

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/biensupernice/krane/internal/store"
)

const boltpath = "./krane.db"
const namespace = "krane_test"

func teardown() { os.Remove(boltpath) }

func TestMain(m *testing.M) {
	store.Connect((boltpath))
	defer store.Client().Disconnect()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestNewEnqueuer(t *testing.T) {
	jobChannel := make(chan Job)

	e := NewEnqueuer(jobChannel)

	assert.NotNil(t, e)
	assert.Equal(t, jobChannel, e.queue)
}

func TestEnqueueNewJobs(t *testing.T) {
	jobQueue := make(chan Job)

	e := NewEnqueuer(jobQueue)

	jobCount := 10
	var jobHandlerCalls int

	// Act
	go func(handler *int) {
		for i := 0; i < jobCount; i++ {
			job := Job{
				ID:         string(i),
				Deployment: namespace,
				Type:       "test",
				Args:       map[string]string{"name": "test"},
				Run: func(args interface{}) error {
					assert.Equal(t, "test", args.(map[string]string)["name"])
					*handler += 1
					return nil
				},
			}

			_, err := e.Enqueue(job)
			assert.Nil(t, err)
			time.Sleep(1 * time.Second)
		}
	}(&jobHandlerCalls)

	// Assert
	for i := 0; i < jobCount; i++ {
		j := <-jobQueue
		j.Run(j.Args)
		assert.NotNil(t, j)
		assert.Equal(t, j.ID, string(i))
		assert.Equal(t, j.Deployment, namespace)
		assert.Equal(t, j.Args.(map[string]string)["name"], "test")
	}

	assert.Equal(t, jobCount, jobHandlerCalls)
}
