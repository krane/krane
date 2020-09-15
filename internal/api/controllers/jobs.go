package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/job"
)

func GetRecentJobs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	daysAgo := query.Get("days_ago")

	//  Defaults to 1 day ago
	if daysAgo == "" {
		daysAgo = "1"
	}
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	recentJobs, err := job.GetJobs(uint(daysAgoNum))
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, recentJobs)
	return
}

func GetJobsByNamespace(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]

	query := r.URL.Query()
	daysAgo := query.Get("days_ago")

	//  defaults to 1 day ago
	if daysAgo == "" {
		daysAgo = "1"
	}
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
	// deployment name
	params := mux.Vars(r)
	namespace := params["namespace"]
	jobID := params["id"]

	j, err := job.GetJobByID(namespace, jobID)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, j)
	return
}
