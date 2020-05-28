package event

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/biensupernice/krane/internal/deployment/spec"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/gorilla/websocket"
)

// Clients currently connected
var Clients = make(map[string][]*websocket.Conn)

// Channels
var eventsChannel = make(chan *Event)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Event message structure for a deployment event
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Spec      spec.Spec `json:"deployment"`
}

// Emit : an event about a deployment
func Emit(msg string, s spec.Spec) {
	event := &Event{
		Timestamp: time.Now(),
		Message:   msg,
		Spec:      s,
	}
	eventsChannel <- event
}

// Subscribe : to deployment events
func Subscribe(client *websocket.Conn, deployment string) {
	Clients[deployment] = append(Clients[deployment], client)
}

// Unsubscribe : from deployment events
func Unsubscribe(client *websocket.Conn, deployment string) {
	for i, c := range Clients[deployment] {
		if c == client {
			// Close client connection
			client.Close()

			// Slice out the client
			Clients[deployment] = append(Clients[deployment][:i], Clients[deployment][i+1:]...)
		}
	}
}

// Echo : broadcast deployment events
func Echo() {
	for {
		event := <-eventsChannel
		eventBytes, err := json.Marshal(event)
		if err != nil {
			logger.Debugf("Unable to transform event - %s", err.Error())
			continue
		}

		// Get all clients listening to the specific deployment
		deployment := event.Spec.Name
		clients := Clients[deployment]

		// send to every client that is currently connected
		for _, client := range clients {
			// Write the message to the events channel for deployments
			err := client.WriteMessage(websocket.TextMessage, eventBytes)
			if err != nil {
				// May also want to remove client ???
				logger.Debugf("Websocket error: %s", err)
				Unsubscribe(client, deployment)
			}
		}
	}
}
