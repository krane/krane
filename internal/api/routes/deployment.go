package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/deployment"
)

// RunDeployment : run a deployment
func RunDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()

	name := params["name"]
	_ = query.Get("tag")

	// Find the deployment
	_, err := deployment.GetDeployment(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	// jobs.StartDeployment(deployment, tag)

	status.HTTPAccepted(w)
	return
}

func CreateDeployment(w http.ResponseWriter, r *http.Request) {
	status.HTTPOk(w, nil)
	return
}

// DeleteDeployment : delete a deployment
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	_, err := deployment.GetDeployment(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	// jobs.DeleteDeployment(d)

	status.HTTPAccepted(w)
	return
}

// StopDeployment : stop a deployment
func StopDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	_, err := deployment.GetDeployment(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	// jobs.StopDeployment(d)

	status.HTTPAccepted(w)
	return
}

// GetDeployment : get a deployment
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, d)
	return
}

// GetDeployments : get all deployments
func GetDeployments(w http.ResponseWriter, r *http.Request) {
	// Find deployments
	deployments, err := deployment.GetDeployments()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, deployments)
	return
}
