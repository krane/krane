package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	httpR "github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/deployment"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// CreateDeployment : create a deployment with a template
func CreateDeployment(c *gin.Context) {
	var deploymentTemplate deployment.Template
	err := c.ShouldBindJSON(&deploymentTemplate)
	if err != nil {
		httpR.BadRequest(c, err.Error())
		return
	}

	// Check if deployment already exist
	d := *deployment.FindTemplate(deploymentTemplate.Name)
	if d != (deployment.Template{}) {
		errMsg := fmt.Sprintf("Deployment %s already exist", d.Name)
		httpR.BadRequest(c, errMsg)
		return
	}

	// Save deployment
	err = deployment.SaveTemplate(&deploymentTemplate)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to save deployment - %s", err.Error())
		httpR.BadRequest(c, errMsg)
		return
	}

	httpR.Created(c, deploymentTemplate)
}

// RunDeployment :
func RunDeployment(c *gin.Context) {
	name := c.Param("name")
	tag := c.DefaultQuery("tag", "latest")

	// Check if deployment exist
	deploymentTemplate := *deployment.FindTemplate(name)
	if deploymentTemplate == (deployment.Template{}) {
		errMsg := fmt.Sprintf("Unable to find deployment %s", name)
		httpR.BadRequest(c, errMsg)
		return
	}

	// Start deployment context
	ctx := context.Background()

	// Start deployment using the deployment template and provided tag
	go deployment.Start(&ctx, deploymentTemplate, tag)

	httpR.Accepted(c)
}

// GetDeployment : get single deployment by name
func GetDeployment(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		errMsg := errors.New("Error getting deployment `name` from params")
		httpR.BadRequest(c, errMsg)
		return
	}

	// Get deployment by name
	d := deployment.FindTemplate(name)

	// compare an empty deployment against the one found in the store
	if *d == (deployment.Template{}) {
		httpR.BadRequest(c, "Unable to find a deployment by that name")
		return
	}

	httpR.Ok(c, &d)
}

// GetDeployments : get all deployments
func GetDeployments(c *gin.Context) { httpR.Ok(c, deployment.FindAllTemplates()) }

// DeleteDeployment : delete deployment and its resources
func DeleteDeployment(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		errMsg := errors.New("Error getting deployment `name` from params")
		httpR.BadRequest(c, errMsg)
		return
	}

	// Get deployment by name
	d := deployment.FindTemplate(name)

	// compare an empty deployment against the one found in the store
	if *d == (deployment.Template{}) {
		httpR.BadRequest(c, "Unable to find a deployment by that name")
		return
	}

	ctx := context.Background()

	// Delete a deployments docker resources
	go deployment.DeleteDockerResources(&ctx, *d)

	// Delete deployment from data store
	store.Remove(store.DeploymentsBucket, name)

	ctx.Done()

	httpR.Accepted(c)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSDeploymentHandler :  handler for incoming websocket requests
func WSDeploymentHandler(c *gin.Context) {
	name := c.Param("name")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Debugf("Error registering client - %s", err.Error())
		deployment.Unsubscribe(ws, name)
		return
	}

	deployment.Subscribe(ws, name)
	logger.Debugf("Registered new client - %s", ws)
}
