package deployment

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/krane/krane/internal/logger"
)

type EventEmitter struct {
	Deployment string
	JobID      string
	Clients    []*websocket.Conn
}

type Event struct {
	JobID      string    `json:"job_id"`
	Deployment string    `json:"deployment"`
	Type       EventType `json:"type"`
	Message    string    `json:"message"`
}

type EventType string

const (
	DeploymentSetup           EventType = "DEPLOYMENT_SETUP"
	DeploymentHealthCheck     EventType = "DEPLOYMENT_HEALTHCHECK"
	DeploymentCleanup         EventType = "DEPLOYMENT_CLEANUP"
	DeploymentDone            EventType = "DEPLOYMENT_DONE"
	DeploymentPullImage       EventType = "PULL_IMAGE"
	DeploymentCreateContainer EventType = "CREATE_CONTAINER"
	DeploymentStartContainer  EventType = "START_CONTAINER"
	DeploymentError           EventType = "ERROR"
)

var eventClients = make(map[string][]*websocket.Conn)
var eventsMutex = &sync.Mutex{}

func createEventEmitter(deployment string, jobID string) *EventEmitter {
	return &EventEmitter{
		Deployment: deployment,
		JobID:      jobID,
		Clients:    eventClients[deployment],
	}
}

// emit broadcasts an event payload to all clients connected to that deployment.
// In order to allow clients to filter events for specific deployment runs, the job id
// was added into the event payload, the job id is returned when triggering a deployment run.
func (e EventEmitter) emit(eventType EventType, message string) {
	go func(clients []*websocket.Conn, jobID string, deployment string) {
		for _, client := range clients {
			bytes, _ := json.Marshal(Event{
				JobID:      jobID,
				Deployment: e.Deployment,
				Type:       eventType,
				Message:    message,
			})

			eventsMutex.Lock()
			if err := client.WriteMessage(websocket.TextMessage, bytes); err != nil {
				// this will log when a client has disconnected at which point the
				// connection is not valid causing a write error. This should not
				// affect other clients or streaming logs in general.
				logger.Debugf("client %v disconnected: %v", client.RemoteAddr(), err)
				UnSubscribeFromDeploymentEvents(client, deployment)
			}
			eventsMutex.Unlock()
		}
	}(e.Clients, e.JobID, e.Deployment)
}

// emitStream broadcast a stream of data to all clients connected to the deployment.
// A stream could be the data when pulling an image, reading container logs etc... where an io.Reader is returned
func (e EventEmitter) emitStream(eventType EventType, reader io.Reader) {
	buffReader := bufio.NewReader(reader)
	for {
		bytes, _, err := buffReader.ReadLine()
		if err != nil {
			return
		}

		data, _ := json.Marshal(Event{
			JobID:      e.JobID,
			Deployment: e.Deployment,
			Type:       eventType,
			Message:    string(bytes),
		})
		for _, client := range e.Clients {
			eventsMutex.Lock()
			if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
				// this will log when a client has disconnected at which point the
				// connection is not valid causing a write error. This should not
				// affect other clients or streaming logs in general.
				logger.Debugf("client %v disconnected: %v", client.RemoteAddr(), err)
				UnSubscribeFromDeploymentEvents(client, e.Deployment)
			}
			eventsMutex.Unlock()
		}
	}
}

func (e EventEmitter) closeStream(eventType EventType, message string) {
	for _, client := range e.Clients {
		data, _ := json.Marshal(Event{
			JobID:      e.JobID,
			Deployment: e.Deployment,
			Type:       eventType,
			Message:    message,
		})
		eventsMutex.Lock()
		if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
			// this will log when a client has disconnected at which point the
			// connection is not valid causing a write error. This should not
			// affect other clients or streaming logs in general.
			logger.Debugf("client %v disconnected: %v", client.RemoteAddr(), err)
		}
		UnSubscribeFromDeploymentEvents(client, e.Deployment)
		eventsMutex.Unlock()
	}
}

// SubscribeToDeploymentEvents allows clients to subscribes to a particular deployments events
func SubscribeToDeploymentEvents(client *websocket.Conn, deployment string) {
	eventClients[deployment] = append(eventClients[deployment], client)

	// This will read indefinitely until the client closes
	// the connection ensuring we cleanup up dead connections.
	// Clients should invoke `ws.close()` so that the server
	// can properly unsubscribe from deployment events.
	go func(client *websocket.Conn, deployment string) {
		for {
			if _, _, err := client.NextReader(); err != nil {
				UnSubscribeFromDeploymentEvents(client, deployment)
				break
			}
		}
	}(client, deployment)
}

// UnSubscribeFromDeploymentEvents unsubscribes a client from deployment events
func UnSubscribeFromDeploymentEvents(client *websocket.Conn, deployment string) {
	for i, c := range eventClients[deployment] {
		if c == client {
			if err := client.Close(); err != nil {
				logger.Warnf("unable to properly close client connection %v", err)
			}
			eventClients[deployment] = append(eventClients[deployment][:i], eventClients[deployment][i+1:]...)
		}
	}
}
