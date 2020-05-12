package server

import (
	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

type DeploymentResponse struct {
	Success bool              `json:"success"`
	Spec    docker.DeploySpec `json:"spec"`
}

// DeployAppHandler : deploy an app with a deployment spec
func DeployAppHandler(c *gin.Context) {
	var spec docker.DeploySpec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	err = docker.Deploy(spec)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	http.Ok(c, &DeploymentResponse{
		Success: true,
		Spec:    spec,
	})
}
