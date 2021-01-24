package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/deployment"
	"github.com/krane/krane/internal/session"
)

// WSUpgrader upgrades HTTP connections to WebSocket connections
var WSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		s := r.Context().Value("session").(session.Session)
		if !s.IsValid() {
			return false
		}
		return true
	},
}

// GetDeployment returns a deployments and it configuration, containers and recent activity
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

	d, err := deployment.GetDeployment(deploymentName)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, d)
	return
}

// GetAllDeployments returns a list of deployments with their configurations, containers and recent activity
func GetAllDeployments(w http.ResponseWriter, _ *http.Request) {
	deployments, err := deployment.GetAllDeployments()
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, deployments)
	return
}

// CreateOrUpdateDeployment saves a deployment configuration
func CreateOrUpdateDeployment(w http.ResponseWriter, r *http.Request) {
	var config deployment.Config

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		response.HTTPBad(w, err)
		return
	}

	if err := deployment.SaveConfig(config); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, config)
	return
}

// DeleteDeployment deletes a deployments container resources and configuration
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

// RunDeployment triggers a deployment run creating container resources
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

// GetDeploymentContainers returns all containers for a deployment
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

// ReadContainerLogs upgrades the incoming http connection into a websocket connection to stream container logs
func ReadContainerLogs(w http.ResponseWriter, r *http.Request) {
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

// SubscribeToDeploymentEvents upgrades the incoming http connection into a websocket connection to stream deployment events
func SubscribeToDeploymentEvents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deploymentName := params["deployment"]

	connection, err := WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	deployment.SubscribeToDeploymentEvents(connection, deploymentName)
	return
}
