package deployment

import (
	"github.com/gorilla/websocket"

	"github.com/krane/krane/internal/docker"
	"github.com/krane/krane/internal/logger"
)

// SubscribeToDeploymentLogs streams deployment logs to a websocket client
func SubscribeToDeploymentLogs(client *websocket.Conn, deployment string) {
	data := make(chan []byte)
	done := make(chan bool)

	containers, err := GetContainersByDeployment(deployment)
	if err != nil {
		logger.Warnf("unable to get containers for deployment %s, %v", deployment, err)
		if err := client.Close(); err != nil {
			logger.Warnf("error closing client connection, %v", err)
			return
		}
		return
	}

	for _, container := range containers {
		if err := docker.GetClient().StreamContainerLogs(container.ID, data, done); err != nil {
			logger.Warnf("error grabbing container reader, %v", err)
			if err := client.Close(); err != nil {
				logger.Warnf("error closing client connection, %v", err)
				return
			}
			return
		}
	}

	for {
		select {
		case bytes := <-data:
			if err := client.WriteMessage(websocket.TextMessage, bytes); err != nil {
				// this will log when a client has disconnected at which point the
				// connection is not valid causing a write error. This should not
				// affect other clients or streaming logs in general.
				logger.Debugf("client %v disconnected", client.LocalAddr())
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

// SubscribeToContainerLogs streams container logs to a websocket client
func SubscribeToContainerLogs(client *websocket.Conn, containerID string) {
	data := make(chan []byte)
	done := make(chan bool)

	if err := docker.GetClient().StreamContainerLogs(containerID, data, done); err != nil {
		logger.Warnf("error grabbing container reader, %v", err)
		if err := client.Close(); err != nil {
			logger.Warnf("error closing client connection, %v", err)
			return
		}
		return
	}

	for {
		select {
		case bytes := <-data:
			if err := client.WriteMessage(websocket.TextMessage, bytes); err != nil {
				// this will log when a client has disconnected at which point the
				// connection is not valid causing a write error. This should not
				// affect other clients or streaming logs in general.
				logger.Debugf("client %v disconnected", client.LocalAddr())
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
