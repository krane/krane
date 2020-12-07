package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/namespace"
	"github.com/biensupernice/krane/internal/deployment/service"
)

// ApplyDeployment : create or update a deployment
func ApplyDeployment(w http.ResponseWriter, r *http.Request) {
	var cfg config.DeploymentConfig

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

	status.HTTPAccepted(w)
	return
}

// DeleteDeployment : delete a deployment, removing the container resources and configuration
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	cfg, err := config.GetDeploymentConfig(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	if err := service.DeleteDeployment(cfg); err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPAccepted(w)
	return
}

// GetContainers : gets all containers for a deployment
func GetContainers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	if name == "" {
		status.HTTPBad(w, errors.New("deployment name required"))
		return
	}

	if !namespace.Exist(name) {
		status.HTTPBad(w, errors.New("deployment does not exist"))
		return
	}

	containers, err := container.GetContainersByDeployment(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, containers)
	return
}

// GetDeploymentConfig : get a deployment configuration
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	deployment, err := config.GetDeploymentConfig(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, deployment)
	return
}

// GetDeploymentConfig : get all deployments
func GetAllDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := config.GetAllDeploymentConfigurations()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, deployments)
	return
}
