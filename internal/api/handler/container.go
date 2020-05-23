package handler

import (
	"context"
	"fmt"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/gin-gonic/gin"
)

// ListContainers : list all containers
func ListContainers(c *gin.Context) {
	ctx := context.Background()
	containers, err := docker.ListContainers(&ctx)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	http.Ok(c, containers)
}

// StopContainer : stop docker container
func StopContainer(c *gin.Context) {
	containerID := c.Param("containerID")

	if containerID == "" {
		http.BadRequest(c, "Container ID required")
		return
	}

	ctx := context.Background()

	err := docker.StopContainer(&ctx, containerID)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to stop container %s", containerID)
		http.BadRequest(c, errMsg)
		return
	}

	ctx.Done()

	msg := fmt.Sprintf("Container %s stopped", containerID)
	logger.Debug(msg)

	http.Ok(c, map[string]string{"message": msg})
}

// StartContainer : start docker container
func StartContainer(c *gin.Context) {
	containerID := c.Param("containerID")

	if containerID == "" {
		http.BadRequest(c, "Container ID required")
		return
	}

	ctx := context.Background()

	err := docker.StartContainer(&ctx, containerID)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to start container %s - %s", containerID, err.Error())
		http.BadRequest(c, errMsg)
		return
	}

	ctx.Done()

	msg := fmt.Sprintf("Container %s started", containerID)
	logger.Debug(msg)

	http.Ok(c, map[string]string{"message": msg})
}
