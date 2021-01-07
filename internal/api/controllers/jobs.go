package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/deployment"
	"github.com/krane/krane/internal/utils"
)

// GetRecentJobs returns all deployment jobs within a date range (default is 1d ago)
func GetRecentJobs(w http.ResponseWriter, r *http.Request) {
	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "1")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	jobs, err := deployment.GetJobs(uint(daysAgoNum))
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, jobs)
	return
}

// GetJobsByDeployment returns all jobs within a date range (default is 1d ago)
func GetJobsByDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "1")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
		return
	}

	jobs, err := deployment.GetJobsByDeployment(deploymentName, uint(daysAgoNum))
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, jobs)
	return
}

// GetJobByID returns a job by id
func GetJobByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]
	jobID := params["id"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if jobID == "" {
		response.HTTPBad(w, errors.New("job id not provided"))
		return
	}

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
		return
	}

	daysAgo := utils.QueryParamOrDefault(r, "days_ago", "365")
	daysAgoNum, _ := strconv.Atoi(daysAgo)

	j, err := deployment.GetJobByID(deploymentName, jobID, uint(daysAgoNum))
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, j)
	return
}
