package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/service"
)

// SaveDeployment: save a deployment creating container resources
func SaveDeployment(w http.ResponseWriter, r *http.Request) {
	var cfg config.Config

	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		status.HTTPBad(w, err)
		return
	}

	if err := cfg.Save(); err != nil {
		status.HTTPBad(w, err)
		return
	}

	if err := service.StartDeployment(cfg); err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, cfg)
	return
}

// DeleteDeployment: delete a deployment removing the container resources and configuration
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	d, err := service.GetDeploymentByName(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	if err := service.DeleteDeployment(d); err != nil {
		status.HTTPBad(w, err)
		return
	}

	// TODO: remove configuration
	
	status.HTTPOk(w, nil)
	return
}

// GetDeployment: get deployment by name
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	d, err := service.GetDeploymentByName(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, d)
	return
}

// GetDeployment: get all deployments
func GetAllDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := service.GetAllDeployments()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, deployments)
	return
}
