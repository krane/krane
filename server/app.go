package server

import (
	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/http"
	"github.com/gin-gonic/gin"
)

// DeployApp blah
func DeployApp(c *gin.Context) {
	var spec docker.DeploySpec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	res, err := docker.Deploy(spec)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	http.Ok(c, res)
}
