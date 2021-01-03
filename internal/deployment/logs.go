package deployment

import (
	"github.com/gorilla/websocket"

	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/logger"
)

// ReadContainerLogs streams container logs to a websocket client
func ReadContainerLogs(client *websocket.Conn, container string) {
	data := make(chan []byte)
	done := make(chan bool)

	if err := docker.GetClient().StreamContainerLogs(container, data, done); err != nil {
		logger.Warnf("error grabbing container reader, %v", err)
		return
	}

	for {
		select {
		case bytes := <-data:
			if err := client.WriteMessage(websocket.TextMessage, bytes); err != nil {
				// this will log when a client has disconnected at which point the
				// connection is not valid causing a write error. This should not
				// effect other clients or streaming logs in general.
				logger.Debugf("error writing to client socket, %v", err)
				return
			}
		case <-done:
			if err := client.Close(); err != nil {
				logger.Warnf("error closing client connection when unsubscribing from container logs, %v", err)
				return
			}
		}
	}
}
