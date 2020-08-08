package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/service/deployment"
	"github.com/biensupernice/krane/internal/service/jobs"
	"github.com/biensupernice/krane/internal/service/spec"
	"github.com/biensupernice/krane/internal/storage"
)

func RunDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()

	name := params["name"]
	tag := query.Get("tag")

	// Find the deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Create a start deployment job
	jobs.StartDeployment(d, tag)

	utils.HTTPAccepted(w)
	return
}

func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Create a delete deployment job
	jobs.DeleteDeployment(d)

	utils.HTTPAccepted(w)
	return
}

func StopDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Find the deployment
	d, err := deployment.GetDeployment(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Create a stop deployment job
	jobs.StopDeployment(d)

	utils.HTTPAccepted(w)
	return
}

func CreateSpec(w http.ResponseWriter, r *http.Request) {
	var s spec.Spec
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Create spec
	err = s.CreateSpec()
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPCreated(w)
	return
}

func UpdateSpec(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	var s spec.Spec
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Verify the Spec exist
	data, err := storage.Get(spec.Collection, name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	if data == nil {
		utils.HTTPBad(w, errors.New("Deployment with that name not found"))
		return
	}

	// Update spec
	err = s.UpdateSpec(name)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPCreated(w)
	return
}

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

// GetRunningJobs : that are queue'd up
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

func UpdateDeploymentAlias(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	type Body struct {
		Alias string `json:"alias" binding:"required"`
	}

	// Decode request
	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Find deployment
	d, err := deployment.GetDeployment(deploymentName)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Create a update alias job
	jobs.UpdateDeploymentAlias(d, body.Alias)

	utils.HTTPAccepted(w)
	return
}

func DeleteDeploymentAlias(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["name"]

	// Find deployment
	d, err := deployment.GetDeployment(deploymentName)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Create a delete alias job
	jobs.DeleteDeploymentAlias(d)

	utils.HTTPAccepted(w)
	return
}
