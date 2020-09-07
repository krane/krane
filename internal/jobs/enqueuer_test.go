package job

import (
	"os"
	"testing"

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
	// store := store.Instance()
	// jobQueue := make(chan Job)
	//
	// e := NewEnqueuer(store, jobQueue)
	// e.WithHandler("deploy", mockJobHandler)

	// Act
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		e.Enqueue("deploy", map[string]interface{}{"id": i})
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	// Assert
	// go func() {
	// for i := 0; i < 10; i++ {
	// 	j := <-jobQueue
	// 	assert.NotNil(t, j)
	// 	assert.NotNil(t, j.Args["id"])
	// }
	// }()
}

func TestNewEnqueuer(t *testing.T) {
	store := store.Instance()
	jobChannel := make(chan Job)

	e := NewEnqueuer(store, jobChannel)
	e.WithHandler("deploy", mockJobHandler)
	e.WithHandler("delete", mockJobHandler)

	assert.NotNil(t, e)
	assert.Equal(t, store, e.store)
	assert.Equal(t, jobChannel, e.jobQueue)
	assert.NotNil(t, e.Handlers["deploy"])
	assert.NotNil(t, e.Handlers["delete"])
	assert.Nil(t, e.Handlers["update"])
}

func TestErrorThrownWhenNoJobHandlerRegistered(t *testing.T) {
	e := NewEnqueuer(nil, nil)
	_, err := e.Enqueue("deploy", nil)

	assert.NotNil(t, err)
	assert.Error(t, err, "unable to queue job, unknown handler")
}
