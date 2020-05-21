package handler

import (
	"errors"
	"log"

	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/deployment"
	"github.com/gin-gonic/gin"
)

// CreateDeployment : using deployment spec
func CreateDeployment(c *gin.Context) {
	var t deployment.Template
	err := c.ShouldBindJSON(&t)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Compare with a zero value composite literal because all fields are comparable
	d := *deployment.FindDeployment(t.Name)
	if d != (deployment.Template{}) {
		http.BadRequest(c, "Deployment with that name already exists")
		return
	}

	deployment.SaveDeployment(&t)

	// Start new deployment thread
	go deployment.Start2(t)

	http.Accepted(c)
}

// GetDeployment : get single deployment by name
func GetDeployment(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		log.Println("pipi")
		http.BadRequest(c, errors.New("Error getting deployment `name` from params"))
		return
	}

	// Get deployment
	d := deployment.FindDeployment(name)

	// Check if deployment was not found
	if *d == (deployment.Template{}) {
		http.Ok(c, nil)
		return
	}

	http.Ok(c, &d)
}

// GetDeployments : get all deployments
func GetDeployments(c *gin.Context) {
	// Get deployments
	deployments := deployment.FindAllDeployments()

	http.Ok(c, deployments)
}
