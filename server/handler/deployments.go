package handler

import (
	"encoding/json"
	"log"

	"github.com/biensupernice/krane/data"
	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

// CreateDeployment : using deployment spec
func CreateDeployment(c *gin.Context) {
	var deployment docker.Deployment
	err := c.ShouldBindJSON(&deployment)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Queue up the deployment
	go docker.QueueDeployment(deployment)

	http.Accepted(c)
}

// GetDeployments : get all deployments
func GetDeployments(c *gin.Context) {
	// Get deployments
	deploymentData := data.GetAll(data.DeploymentsBucket)

	var deployments []docker.Deployment
	for v := 0; v < len(deploymentData); v++ {
		var d docker.Deployment
		err := json.Unmarshal(*deploymentData[v], &d)
		if err != nil {
			log.Printf("Unable to parse deployment [skipping] - %s", err.Error())
			continue
		}
		deployments = append(deployments, d)
	}

	http.Ok(c, deployments)
}
