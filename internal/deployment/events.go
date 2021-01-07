package deployment

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/gorilla/websocket"

	"github.com/krane/krane/internal/logger"
)

type EventEmitter struct {
	Deployment string
	JobID      string
	Phase      Phase
	Clients    []*websocket.Conn
}

type Event struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
	Phase   Phase  `json:"phase"`
}

var eventClients = make(map[string][]*websocket.Conn)

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
func (e EventEmitter) emit(message string) {
	go func(clients []*websocket.Conn, jobID string, deployment string, phase Phase) {
		for _, client := range clients {
			bytes, _ := json.Marshal(Event{
				JobID:   jobID,
				Message: message,
				Phase:   phase,
			})
			if err := client.WriteMessage(websocket.TextMessage, bytes); err != nil {
				// this will log when a client has disconnected at which point the
				// connection is not valid causing a write error. This should not
				// affect other clients or streaming logs in general.
				logger.Debugf("client %v disconnected", client.RemoteAddr())
				UnSubscribeFromDeploymentEvents(client, deployment)
				return
			}
		}
	}(e.Clients, e.JobID, e.Deployment, e.Phase)
}

func (e EventEmitter) emitS(reader io.Reader) {
	buffReader := bufio.NewReader(reader)
	go func() {
		for {
			bytes, _, err := buffReader.ReadLine()
			if err != nil {
				return
			}

			data, _ := json.Marshal(Event{
				JobID:   e.JobID,
				Message: string(bytes),
				Phase:   e.Phase,
			})
			for _, client := range e.Clients {
				if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
					// this will log when a client has disconnected at which point the
					// connection is not valid causing a write error. This should not
					// affect other clients or streaming logs in general.
					logger.Debugf("client %v disconnected", client.RemoteAddr())
					UnSubscribeFromDeploymentEvents(client, e.Deployment)
					return
				}
			}
		}
	}()
}

// SubscribeToDeploymentEvents allows clients to subscribes to a particular deployments events
func SubscribeToDeploymentEvents(client *websocket.Conn, deployment string) {
	eventClients[deployment] = append(eventClients[deployment], client)
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
