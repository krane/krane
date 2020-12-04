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

func Connect() {
	once.Do(func() { newClientFromEnv() })
}

func newClientFromEnv() {
	logger.Info("Connecting to Docker client")

	envClient, err := client.NewEnvClient()
	if err != nil {
		logger.Fatalf("Failed creating Docker client %s", err.Error())
		return
	}

	instance = &Client{envClient}

	if err := createDockerNetwork(); err != nil {
		logger.Fatalf("Failed creating Docker network %s", err.Error())
		return
	}

	return
}

func createDockerNetwork() error {
	logger.Debug("Creating Krane Docker network")

	ctx := context.Background()
	defer ctx.Done()

	_, err := instance.CreateBridgeNetwork(&ctx, KraneNetworkName)
	if err != nil {
		return err
	}
	return nil
}
