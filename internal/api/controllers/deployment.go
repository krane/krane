package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/service"
)

// applyDeployment : create or update a deployment
func ApplyDeployment(w http.ResponseWriter, r *http.Request) {
	var cfg config.DeploymentConfig

	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		response.HTTPBad(w, err)
		return
	}

	if err := cfg.Save(); err != nil {
		response.HTTPBad(w, err)
		return
	}

	if err := service.StartDeployment(cfg); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

// DeleteDeployment : delete a deployment
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	cfg, err := config.GetDeploymentConfig(name)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	if err := service.DeleteDeployment(cfg); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

// StopDeployment : stop all containers for a deployment
func StopDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	cfg, err := config.GetDeploymentConfig(name)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	if err := service.StopDeployment(cfg); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

// GetDeploymentConfig : get a deployments configuration
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	deployment, err := config.GetDeploymentConfig(name)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, deployment)
	return
}

// GetDeploymentConfig : get all deployments
func GetAllDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := config.GetAllDeploymentConfigurations()
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, deployments)
	return
}
