package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/deployment/service"
)

var config = kconfig.Kconfig{
	Name:  "krane-proxy",
	Image: "biensupernice/proxy",
	Scale: 1,
	Volumes: map[string]string{
		"/var/run/docker.sock": "/var/run/docker.sock",
	},
	Ports: map[string]string{
		"80":   "80",
		"443":  "443",
		"8080": "8080",
	},
}

// CreateProxyIfNotExist : create the network proxy is it does not already exist
func CreateProxyIfNotExist() {
	containers, err := container.GetContainersByNamespace(config.Name)
	if err != nil {
		panic(fmt.Sprintf("Unable to create network proxy, %v", err))
	}

	// create the proxy if no containers are currently up
	if len(containers) == 0 {
		err := createProxy()
		if err != nil {
			// If we cant create the proxy, exit the program
			panic(fmt.Sprintf("Unable to create network proxy, %v", err))
			return
		}
		return
	}

	// create the proxy if no containers are in a running state
	for _, c := range containers {
		if !c.State.Running {
			err := createProxy()
			if err != nil {
				// If we cant create the proxy, exit the program
				panic(fmt.Sprintf("Unable to create network proxy, %v", err))
				return
			}
			return
		}
	}
	logrus.Debug("Network proxy running...")
}

func createProxy() error {
	err := config.Save()
	if err != nil {
		return err
	}

	err = service.StartDeployment(config)
	if err != nil {
		return err
	}
	logrus.Debug("Network proxy deployed")
	return nil
}
