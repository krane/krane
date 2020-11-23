package job

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/store"
)

const boltpath = "./krane.db"
const namespace = "krane_test"

func teardown() { os.Remove(boltpath) }

func TestMain(m *testing.M) {
	store.NewInstance((boltpath))
	defer store.Instance().Shutdown()

	// Create deployment (namespace)
	deployment := kconfig.Kconfig{Name: namespace}
	bytes, _ := deployment.Serialize()
	store.Instance().Put(constants.DeploymentsCollectionName, deployment.Name, bytes)

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestNewEnqueuer(t *testing.T) {
	store := store.Instance()
	jobChannel := make(chan Job)

	e := NewEnqueuer(store, jobChannel)

	assert.NotNil(t, e)
	assert.Equal(t, store, e.store)
	assert.Equal(t, jobChannel, e.queue)
}

func TestEnqueueNewJobs(t *testing.T) {
	store := store.Instance()
	jobQueue := make(chan Job)

	e := NewEnqueuer(store, jobQueue)

	jobCount := 10
	var jobHandlerCalls int

	// Act
	go func(handler *int) {
		for i := 0; i < jobCount; i++ {
			job := Job{
				ID:        string(i),
				Namespace: namespace,
				Type:      "test",
				Args:      Args{"name": "test"},
				Run: func(args Args) error {
					logrus.Printf("Job handler called, %v", args)
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
		assert.Equal(t, j.Namespace, namespace)
		assert.Equal(t, j.Args["name"], "test")
	}

	assert.Equal(t, jobCount, jobHandlerCalls)
}
