package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/utils"
)

func GetRecentJobs(w http.ResponseWriter, r *http.Request) {
	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "1")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	job, err := job.GetJobs(uint(daysAgoNum))
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, job)
	return
}

func GetJobsByNamespace(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]

	if namespace == "" {
		status.HTTPBad(w, errors.New("namespace not provided"))
		return
	}

	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "1")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	jobs, err := job.GetJobsByNamespace(namespace, uint(daysAgoNum))
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, jobs)
	return
}

func GetJobByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	jobID := params["id"]

	if namespace == "" {
		status.HTTPBad(w, errors.New("namespace not provided"))
		return
	}

	if jobID == "" {
		status.HTTPBad(w, errors.New("job id not provided"))
		return
	}

	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "365")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	j, err := job.GetJobByID(namespace, jobID, uint((daysAgoNum)))
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, j)
	return
}
