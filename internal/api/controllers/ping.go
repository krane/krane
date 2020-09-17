package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

func PingController(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	message := params["message"]

	newJob := job.Job{
		ID:          utils.MakeIdentifier(),
		Namespace:   namespace,
		Type:        job.ContainerCreate,
		Args:        map[string]interface{}{"message": message},
		RetryPolicy: 3,
		Run:         ping,
	}

	enqueuer := job.NewEnqueuer(store.Instance(), job.GetJobQueue())
	go func() {
		_, err := enqueuer.Enqueue(newJob)
		if err != nil {
			logrus.Errorf(err.Error())
		}
	}()

	status.HTTPOk(w, "Got message")
	return
}

func ping(args job.Args) error {
	logrus.Infof("Got message %s", args["message"])
	return fmt.Errorf("Test error with args %s", args["message"])
}
