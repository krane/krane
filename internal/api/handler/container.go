package handler

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/biensupernice/krane/internal/channel"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/gin-contrib/sse"
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

// ContainerEvents : stream container event, logs, errors using server-sent-events(Sse)
func ContainerEvents(c *gin.Context) {
	cID := c.Param("containerID")
	if cID == "" {
		http.BadRequest(c, "container id not provided")
		return
	}

	// Start context
	ctx := context.Background()

	// Channel to stream
	containerEvents := make(chan string)
	done := make(chan bool)

	ioreader, err := docker.ReadContainerLogs(&ctx, cID)
	if err != nil {
		logger.Debugf("error reader: %s", err.Error())
		done <- true
		return
	}

	// Stream container events from the reader return from `ReadContainerLogs` to a streaming channel
	// That gin can use to serve to the client a stream using server-sent-events (sse)
	channel.Stream(&ioreader, containerEvents, done)

	msgCount := 0
	c.Stream(func(w io.Writer) bool {
		for {
			select {
			case <-done:
				// when deadline is reached, send 'end' event
				c.SSEvent("end", "end")
				return false
			case event := <-containerEvents:
				// send events to client
				c.Render(-1, sse.Event{
					Id:    strconv.Itoa(msgCount), // Current msg # is id
					Event: "event",
					Data:  event,
				})
				msgCount++
				return true
			}
		}
	})
}
