package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/deployment"
)

// GetDeployment returns the configuration for a single deployment
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	d, err := deployment.GetDeploymentConfig(deploymentName)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, d)
	return
}

// GetAllDeployments returns all deployments and their configurations
func GetAllDeployments(w http.ResponseWriter, _ *http.Request) {
	deployments, err := deployment.GetAllDeploymentConfigs()
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, deployments)
	return
}

// SaveDeployment : create or update a deployment
func SaveDeployment(w http.ResponseWriter, r *http.Request) {
	var config deployment.Config

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		response.HTTPBad(w, err)
		return
	}

	if err := deployment.Save(config); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, config)
	return
}

// DeleteDeployment deletes a deployments containers and configuration
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if err := deployment.Delete(deploymentName); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

func RunDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if err := deployment.Run(deploymentName); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

// GetDeploymentContainers returns all containers for a deployments
func GetDeploymentContainers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, errors.New("deployment does not exist"))
		return
	}

	containers, err := deployment.GetContainersByDeployment(deploymentName)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, containers)
	return
}
