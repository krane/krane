package handler

import (
	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

// DeploymentResponse : from deploying an app
type DeploymentResponse struct {
	ContainerID string `json:"container_id"`
}

// CreateDeployment : using deployment spec
func CreateDeployment(c *gin.Context) {
	var spec docker.DeploySpec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	containerID, err := docker.Deploy(spec)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	http.Ok(c, &DeploymentResponse{ContainerID: containerID})
}

// GetDeployments : get all deployments
func GetDeployments(c *gin.Context) {
	http.Ok(c, "Not yet implemented")
}
