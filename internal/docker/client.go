package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*client.Client
}

var instance *Client
var once sync.Once

// GetClient : get docker client
func GetClient() *Client { return instance }

// Init : starts docker client
func NewClient() *Client {
	if instance != nil {
		return instance
	}

	logrus.Info("Connecting to Docker client...")
	once.Do(func() {
		envClient,
		err := client.NewEnvClient()
		if err != nil {
			logrus.Fatalf("Failed connecting to Docker envClient on host machine %s", err.Error())
			return
		}

		instance = &Client{envClient}

		ctx := context.Background()

		// Create krane docker network
		logrus.Info("Creating Krane Docker network...")
		_, err = instance.CreateBridgeNetwork(&ctx, KraneNetworkName)
		if err != nil {
			logrus.Fatalf("Failed to create Krane Docker network %s", err.Error())
		}

		ctx.Done()
	})

	return instance
}
