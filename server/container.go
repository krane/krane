package server

import (
	"context"
	"fmt"
	"log"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/http"
	"github.com/gin-gonic/gin"
)

// StopContainerHandler blah
func StopContainerHandler(c *gin.Context) {
	containerID := c.Param("containerID")

	if containerID == "" {
		http.BadRequest(c, "Container ID required")
		return
	}

	client, err := docker.NewClient()
	if err != nil {
		msg, _ := fmt.Printf("Unable to start docker client %s", err.Error())
		http.BadRequest(c, msg)
		return
	}
	ctx := context.Background()

	err = docker.StopContainer(&ctx, client, containerID)
	if err != nil {
		msg := fmt.Sprintf("Unable to stop container %s", containerID)
		http.BadRequest(c, msg)
		return
	}

	msg := fmt.Sprintf("Container %s stopped", containerID)
	log.Printf(msg)

	http.Ok(c, msg)
}

// StartContainerHandler blah
func StartContainerHandler(c *gin.Context) {
	containerID := c.Param("containerID")

	if containerID == "" {
		http.BadRequest(c, "Container ID required")
		return
	}

	client, err := docker.NewClient()
	if err != nil {
		msg, _ := fmt.Printf("Unable to start docker client %s", err.Error())
		http.BadRequest(c, msg)
		return
	}
	ctx := context.Background()

	err = docker.StartContainer(&ctx, client, containerID)
	if err != nil {
		msg := fmt.Sprintf("Unable to start container %s - %s", containerID, err.Error())
		http.BadRequest(c, msg)
		return
	}

	msg := fmt.Sprintf("Container %s started", containerID)
	log.Printf(msg)

	http.Ok(c, msg)
}
