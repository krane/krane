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

// SaveDeployment: save a deployment creating container resources
func SaveDeployment(w http.ResponseWriter, r *http.Request) {
	var cfg config.Kconfig

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

// DeleteDeployment: delete a deployment removing the container resources and configuration
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	cfg, err := config.GetConfigByDeploymentByName(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	if err := service.DeleteDeployment(cfg); err != nil {
		status.HTTPBad(w, err)
		return
	}

	// TODO: remove configuration

	status.HTTPAccepted(w)
	return
}

// GetDeploymentContainers : gets all the containers and current state
func GetDeploymentContainers(w http.ResponseWriter, r *http.Request) {
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

	containers, err := container.GetKontainersByNamespace(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, containers)
	return
}

// GetDeployment: get deployment by name
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	cfg, err := config.GetConfigByDeploymentByName(name)
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, cfg)
	return
}

// GetDeployment: get all deployments
func GetAllDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := config.GetAllDeploymentConfigs()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, deployments)
	return
}
