package main

import (
	"os"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/deployment/service"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/utils"
)

var deployment = config.DeploymentConfig{
	Name:    "krane-proxy",
	Image:   "biensupernice/proxy",
	Scale:   1,
	Secured: utils.BoolEnv(constants.EnvProxyDashboardSecure),
	Alias:   []string{os.Getenv(constants.EnvProxyDashboardAlias)},
	Env: map[string]string{
		constants.EnvLetsEncryptEmail: os.Getenv(constants.EnvLetsEncryptEmail),
	},
	Volumes: map[string]string{
		"/var/run/docker.sock": "/var/run/docker.sock",
	},
	Ports: map[string]string{
		"80":   "80",
		"443":  "443",
		"8080": "8080",
	},
}

// EnsureNetworkProxy : ensures the network proxy is up and in a running state
func EnsureNetworkProxy() {
	isEnabled := utils.BoolEnv(constants.EnvProxyEnabled)
	if !isEnabled {
		logger.Info("Network proxy not enabled")
		return
	}

	isSecured := utils.BoolEnv(constants.EnvProxyDashboardSecure)
	if isSecured && os.Getenv(constants.EnvLetsEncryptEmail) == "" {
		logger.Fatalf("Missing required environment variable %s when running in SECURE mode", constants.EnvLetsEncryptEmail)
	}

	// get containers (if any) for the proxy deployment
	containers, err := container.GetKraneContainersByDeployment(deployment.Name)
	if err != nil {
		logger.Fatalf("Unable to create network proxy, %v", err)
	}

	// create the proxy if no containers are currently up
	if len(containers) == 0 {
		err := createProxy()
		if err != nil {
			// If we cant create the proxy, exit the program
			logger.Fatalf("Unable to create network proxy, %v", err)
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
				logger.Fatalf("Unable to create network proxy, %v", err)
				return
			}
			return
		}
	}

	logger.Debug("Network proxy running")
}

func createProxy() error {
	if err := deployment.Save(); err != nil {
		return err
	}

	if err := service.StartDeployment(deployment); err != nil {
		return err
	}

	logger.Debug("Network proxy deployment started")
	return nil
}
