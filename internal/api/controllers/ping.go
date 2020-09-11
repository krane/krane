package controllers

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api/status"
	job "github.com/biensupernice/krane/internal/jobs"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

func PingController(w http.ResponseWriter, r *http.Request) {

	newJob := job.Job{
		ID:         utils.MakeIdentifier(),
		Namespace:  "Ping",
		EnqueuedAt: time.Now().Unix(),
		Args:       map[string]interface{}{"message": "pong"},
		Run:        ping,
	}
	jobEnqueuer := job.NewEnqueuer(store.Instance(), job.GetJobQueue())
	go jobEnqueuer.Enqueue(newJob)

	status.HTTPOk(w, "Got message")
	return
}

func ping(args job.Args) error {
	logrus.Info(args["message"])
	time.Sleep(10 * time.Second)
	return nil
}
