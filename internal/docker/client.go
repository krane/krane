package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/client"

	"github.com/biensupernice/krane/internal/logger"
)

type Client struct{ *client.Client }

var once sync.Once
var instance *Client

func GetClient() *Client { return instance }

// Connect : create a docker client
func Connect() {
	once.Do(func() {
		ClientFromEnv()
		EnsureKraneDockerNetwork()
	})
}

// ClientFromEnv : create a docker client based on environment variables
func ClientFromEnv() {
	logger.Info("Connecting to Docker client")

	c, err := client.NewEnvClient()
	if err != nil {
		logger.Fatalf("Failed creating Docker client %s", err.Error())
		return
	}

	instance = &Client{c}

	return
}

// Ping returns true if the Docker client is actively running
func Ping() bool {
	if instance == nil {
		return false
	}

	ping, err := instance.Client.Ping(context.Background())
	if err != nil || ping.APIVersion == "" {
		return false
	}

	return true
}
