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

// EnsureKraneDockerNetwork : ensure the Krane docker network is created
func EnsureKraneDockerNetwork() {
	ctx := context.Background()
	defer ctx.Done()

	_, err := instance.CreateBridgeNetwork(&ctx, KraneNetworkName)
	if err != nil {
		logger.Fatalf("Unable to create Krane network, %v", err)
	}
}
