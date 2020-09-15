package controllers

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

func PingController(w http.ResponseWriter, r *http.Request) {

	newJob := job.Job{
		ID:          utils.MakeIdentifier(),
		Namespace:   "Ping",
		Type:        job.ContainerCreate,
		Args:        map[string]interface{}{"message": "pong"},
		RetryPolicy: 3,
		Run:         ping,
	}
	jobEnqueuer := job.NewEnqueuer(store.Instance(), job.GetJobQueue())
	go jobEnqueuer.Enqueue(newJob)

	status.HTTPOk(w, "Got message")
	return
}

func ping(args job.Args) error {
	time.Sleep(10 * time.Second)
	logrus.Info(args["message"])
	return nil
}
