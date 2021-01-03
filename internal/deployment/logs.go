package deployment

import (
	"bufio"
	"io"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/logger"
)

var clients = make(map[string][]*websocket.Conn)

// SubscribeToContainerLogs subscribes a websocket client to a containers log output
func SubscribeToContainerLogs(client *websocket.Conn, container string) {
	clients[container] = append(clients[container], client)
	go streamLogs(client, container)
}

// UnsubscribeFromContainerLogs removes a websocket client from steaming container logs
func UnsubscribeFromContainerLogs(client *websocket.Conn, container string) {
	for i, c := range clients[container] {
		if c == client {
			if err := client.Close(); err != nil {
				logger.Warnf("error closing client connection when unsubscribing from container logs, %v", err)
			}
			clients[container] = append(clients[container][:i], clients[container][i+1:]...)
		}
	}
}

func streamLogs(client *websocket.Conn, container string) {
	reader, err := docker.GetClient().StreamContainerLogs(container)
	if err != nil {
		logger.Warnf("error streaming logs, %v", err)
	}

	logs := make(chan string)
	done := make(chan bool)

	toChannel(&reader, logs, done)

	for {
		select {
		case data := <-logs:
			if err := client.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
				logger.Debugf("error writing to socket client", err)
			}
		case <-done:
			UnsubscribeFromContainerLogs(client, container)
		}
	}
}

// toChannel streams data from an io.Reader to a channel
func toChannel(in *io.Reader, out chan string, done chan bool) {
	reader := bufio.NewReader(*in)

	var mu sync.RWMutex
	go func() {
		for {
			mu.Lock()

			header := make([]byte, 8)
			_, err := reader.Read(header)
			if err != nil {
				done <- true
				logger.Debugf("error reading container logs header, %v", err)
				return
			}

			data, _, err := reader.ReadLine()
			if err != nil {
				logger.Debugf("error streaming, %v", err)
				done <- true
				return
			}

			out <- string(data)

			mu.Unlock()
		}
	}()
}
