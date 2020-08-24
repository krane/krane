package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/service/deployment"
	"github.com/biensupernice/krane/internal/service/jobs"
)

// RunDeployment : run a deployment
func RunDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()

	name := params["name"]
	tag := query.Get("tag")

	// Find the deployment
	deployment, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	jobs.StartDeployment(deployment, tag)

	utils.HTTPAccepted(w)
	return
}

// DeleteDeployment : delete a deployment
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	jobs.DeleteDeployment(d)

	utils.HTTPAccepted(w)
	return
}

// StopDeployment : stop a deployment
func StopDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	jobs.StopDeployment(d)

	utils.HTTPAccepted(w)
	return
}

// GetDeployment : get a deployment
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, d)
	return
}

// GetDeployments : get all deployments
func GetDeployments(w http.ResponseWriter, r *http.Request) {
	// Find deployments
	deployments, err := deployment.GetDeployments()
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, deployments)
	return
}

// GetRunningJobs : get running jobs that are currently queue'd up
func GetRunningJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := jobs.GetRunningJobs()

	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, jobs)
	utils.HTTPOk(w, nil)
	return
}
