package job

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/biensupernice/krane/internal/store"
)

const boltpath = "./krane.db"

func teardown() { os.Remove(boltpath) }

func TestMain(m *testing.M) {
	store.New((boltpath))
	defer store.Instance().Shutdown()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func mockJobHandler(args Args) error { return nil }

func TestEnqueueNewJobs(t *testing.T) {
	store := store.Instance()
	jobQueue := make(chan Job)

	e := NewEnqueuer(store, jobQueue)

	// Act
	go func() {
		for i := 0; i < 10; i++ {
			e.Enqueue(Job{ID: string(i), Namespace: "krane", Args: Args{"name": "app1"}})
			time.Sleep(1 * time.Second)
		}
	}()

	// Assert
	for i := 0; i < 10; i++ {
		j := <-jobQueue
		assert.NotNil(t, j)
		assert.Equal(t, j.ID, string(i))
		assert.Equal(t, j.Namespace, "krane")
		assert.Equal(t, j.Args["name"], "app1")
	}
}

func TestNewEnqueuer(t *testing.T) {
	store := store.Instance()
	jobChannel := make(chan Job)

	e := NewEnqueuer(store, jobChannel)

	assert.NotNil(t, e)
	assert.Equal(t, store, e.store)
	assert.Equal(t, jobChannel, e.queue)
}
