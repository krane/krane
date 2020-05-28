package handler

import (
	"context"
	"io"
	"strconv"

	"github.com/biensupernice/krane/internal/channel"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

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
