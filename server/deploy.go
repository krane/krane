package server

import (
	"github.com/biensupernice/krane/deploy"
	"github.com/biensupernice/krane/http"
	"github.com/gin-gonic/gin"
)

func DeployHandler(c *gin.Context) {
	var spec deploy.DeploySpec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Set image tag to `latest` if not provided
	if spec.Config.Tag == "" {
		spec.Config.Tag = "latest"
	}

	deploy.Deploy(spec)

	http.Ok(c, spec)
}
