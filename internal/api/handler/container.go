package handler

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"

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

// StreamContainerLogs : follow the container logs through server-sent-events(Sse)
func StreamContainerLogs(c *gin.Context) {
	cID := c.Param("containerID")
	if cID == "" {
		return
	}

	// Start context
	ctx := context.Background()

	chanStream := make(chan string)
	done := make(chan bool)

	ioreader, err := docker.ReadContainerLogs(&ctx, cID)
	if err != nil {
		logger.Debugf("error reader: %s", err.Error())
		done <- true
		return
	}

	reader := bufio.NewReader(ioreader)
	var mu sync.RWMutex
	go func() {
		for {
			mu.Lock()
			// read lines from the reader
			str, _, err := reader.ReadLine()
			if err != nil {
				logger.Debugf("Read Error: %s", err.Error())
				done <- true
				return
			}
			// send the lines to channel
			chanStream <- string(str)
			mu.Unlock()
		}
	}()

	count := 0 // to indicate the message id
	isStreaming := c.Stream(func(w io.Writer) bool {
		for {
			select {
			case <-done:
				// when deadline is reached, send 'end' event
				c.SSEvent("end", "end")
				return false
			case log := <-chanStream:
				// send events to client
				c.Render(-1, sse.Event{
					Id:    strconv.Itoa(count),
					Event: "log",
					Data:  log,
				})
				count++
				return true
			}
		}
	})

	if !isStreaming {
		logger.Debug("stream closed")
	}
}
