package handler

import (
	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

// DeploymentResponse : response from deploying an app
type DeploymentResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

// CreateDeployment : start a deploymet of an app
func CreateDeployment(c *gin.Context) {
	var spec docker.DeploySpec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		resp := &DeploymentResponse{Success: false, Error: err}
		http.BadRequest(c, resp)
		return
	}

	err = docker.Deploy(spec)
	if err != nil {
		resp := &DeploymentResponse{Success: false, Error: err}
		http.BadRequest(c, resp)
		return
	}

	resp := &DeploymentResponse{Success: true}

	http.Ok(c, resp)
}

// GetDeployments : get all deployments
func GetDeployments(c *gin.Context) {
	http.Ok(c, "Middleware activated")
}
