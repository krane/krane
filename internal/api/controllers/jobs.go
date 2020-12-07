package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/utils"
)

// GetRecentJobs : get jobs within a date range (default is 1d ago)
func GetRecentJobs(w http.ResponseWriter, r *http.Request) {
	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "1")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	jobs, err := job.GetJobs(uint(daysAgoNum))
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, jobs)
	return
}

// GetJobsByNamespace : get jobs by deployment namespace
func GetJobsByNamespace(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]

	if namespace == "" {
		response.HTTPBad(w, errors.New("namespace not provided"))
		return
	}

	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "1")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	jobs, err := job.GetJobsByNamespace(namespace, uint(daysAgoNum))
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, jobs)
	return
}

// GetJobByID : get job by id
func GetJobByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	jobID := params["id"]

	if namespace == "" {
		response.HTTPBad(w, errors.New("namespace not provided"))
		return
	}

	if jobID == "" {
		response.HTTPBad(w, errors.New("job id not provided"))
		return
	}

	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "365")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	j, err := job.GetJobByID(namespace, jobID, uint(daysAgoNum))
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, j)
	return
}
