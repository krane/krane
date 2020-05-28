package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/biensupernice/krane/api/response"
	"github.com/biensupernice/krane/internal/deployment"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/event"
	"github.com/biensupernice/krane/internal/deployment/spec"
	"github.com/biensupernice/krane/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// CreateSpec : creates deployment spec
func CreateSpec(c *gin.Context) {
	// Bind request body to Spec
	var spec spec.Spec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		logger.Debugf("Unable to bind spec - %s", err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	// Create spec, if it already exists it will not update it
	spec.SetDefaults()
	err = spec.Create()
	if err != nil {
		logger.Debugf("Unable to create spec %s- %s", spec.Name, err.Error())
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, spec)
}

// DeleteDeployment : delete a deployment spec and its resources
func DeleteDeployment(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		errMsg := errors.New("Error getting deployment `name` from params")
		response.BadRequest(c, errMsg)
		return
	}

	// Find the spec
	s := spec.Find(name)

	if s == (spec.Spec{}) {
		errMsg := fmt.Sprintf("Unable to find deployment %s", name)
		response.BadRequest(c, errMsg)
		return
	}

	ctx := context.Background()

	// Remove deployment resources
	deployment.Remove(&ctx, s)

	ctx.Done()

	// Remove the spec
	s.Delete()

	response.Accepted(c)
}

// RunDeployment :
func RunDeployment(c *gin.Context) {
	name := c.Param("name")
	tag := c.DefaultQuery("tag", "latest")

	// // Check if deployment exist
	s := spec.Find(name)
	if s == (spec.Spec{}) {
		errMsg := fmt.Sprintf("Unable to find deployment %s", name)
		response.BadRequest(c, errMsg)
		return
	}

	// // Start deployment context
	ctx := context.Background()

	// // Start deployment using the deployment template and provided tag
	deployment.Start(&ctx, s, tag)

	ctx.Done()

	response.Accepted(c)
}

// GetDeployment : get single deployment by name
func GetDeployment(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		errMsg := errors.New("Error getting deployment `name` from params")
		response.BadRequest(c, errMsg)
		return
	}

	// Get deployment by name
	// and compare against empty struct
	s := spec.Find(name)
	if s == (spec.Spec{}) {
		response.BadRequest(c, "Unable to find a deployment by that name")
		return
	}

	ctx := context.Background()

	// Get deployment containers
	containers := container.Get(&ctx, name)

	ctx.Done()

	response.Ok(c, &deployment.Deployment{
		Spec:       s,
		Containers: containers,
	})
}

// GetDeployments : get all deployments
func GetDeployments(c *gin.Context) { response.Ok(c, spec.FindAll()) }

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSDeploymentHandler :  handler to handler clients listening to deployment events
func WSDeploymentHandler(c *gin.Context) {
	name := c.Param("name")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Debugf("Error registering client - %s", err.Error())
		event.Unsubscribe(ws, name)
		return
	}

	event.Subscribe(ws, name)
	logger.Debugf("Registered new client")
}
