package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Client struct{ *client.Client }

var instance *Client
var once sync.Once

func GetClient() *Client { return instance }

func ClientFromEnv() *Client {
	if instance != nil {
		return instance
	}

	once.Do(func() { newClientFromEnv() })
	return instance
}

func newClientFromEnv() {
	logrus.Info("Connecting to Docker client...")

	envClient, err := client.NewEnvClient()
	if err != nil {
		logrus.Fatalf("Failed creating Docker client %s", err.Error())
		return
	}

	instance = &Client{envClient}

	if err := createDockerNetwork(); err != nil {
		logrus.Fatalf("Failed creating Docker network %s", err.Error())
		return
	}

	return
}

func createDockerNetwork() error {
	logrus.Debug("Creating Krane Docker network...")

	ctx := context.Background()
	defer ctx.Done()

	_, err := instance.CreateBridgeNetwork(&ctx, KraneNetworkName)
	if err != nil {
		return err
	}
	return nil
}
