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

// DeployApp : deploy an app with a deployment spec
func DeployApp(c *gin.Context) {
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
