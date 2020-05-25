package deployment

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/biensupernice/krane/internal/logger"
	"github.com/gorilla/websocket"
)

// Client :
type Client struct {
	Valid bool `json:"valid"`
}

// Clients currently connected
var Clients = make(map[string][]*websocket.Conn)

// Channels
var eventsChannel = make(chan *Event)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Event message structure for a deployment
type Event struct {
	Timestamp  time.Time `json:"timestamp"`
	Message    string    `json:"message"`
	Deployment Template  `json:"deployment"`
}

// EmitEvent send a message to the events channel for a deployment
func EmitEvent(msg string, t Template) {
	event := &Event{
		Timestamp:  time.Now(),
		Message:    msg,
		Deployment: t,
	}
	eventsChannel <- event
}

// Subscribe to deployment events
func Subscribe(client *websocket.Conn, deployment string) {
	Clients[deployment] = append(Clients[deployment], client)
}

// Unsubscribe client from events channel
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

// EchoEvents listen and broadcast events to the deployment events channel
func EchoEvents() {
	for {
		event := <-eventsChannel
		eventBytes, err := json.Marshal(event)
		if err != nil {
			logger.Debugf("Unable to transform event - %s", err.Error())
			continue
		}

		// Get all clients listening to the specific deployment
		deployment := event.Deployment.Name
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
