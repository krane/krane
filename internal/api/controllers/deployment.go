package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/deployment"
)

// WSUpgrader upgrades HTTP connections to WebSocket connections
var WSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: authZ
		return true
	},
}

// GetDeployment returns the configuration for a single deployment
func GetDeployment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
		return
	}

	deploymentConfig, err := deployment.GetDeploymentConfig(deploymentName)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, deploymentConfig)
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

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
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

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
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
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
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

// StartDeploymentContainers starts all containers (if any) for a deployment
// Note: this does not create any containers, only start already existing ones
func StartDeploymentContainers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
		return
	}

	if err := deployment.StartContainers(deploymentName); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

// StopDeploymentContainers stops all containers (if any) for a deployment
// Note: this does not create any containers, only stops already existing ones
func StopDeploymentContainers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
		return
	}

	if err := deployment.StopContainers(deploymentName); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

// RestartDeploymentContainers re-creates all containers for a deployment
// Note: this is the same as calling /deployments/{deployment} since both re-create container resources
func RestartDeploymentContainers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	if deploymentName == "" {
		response.HTTPBad(w, errors.New("deployment name not provided"))
		return
	}

	if !deployment.Exist(deploymentName) {
		response.HTTPBad(w, fmt.Errorf("deployment %s does not exist", deploymentName))
		return
	}

	if err := deployment.RestartContainers(deploymentName); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPAccepted(w)
	return
}

// StreamContainerLogs opens a websocket connection to stream
// the logs for a container. It upgrades the incoming http connection
// into a websocket connection and keeps it open as long as the client is listening.
func StreamContainerLogs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	container := params["container"]

	connection, err := WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	deployment.ReadContainerLogs(connection, container)
	return
}
