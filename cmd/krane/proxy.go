package main

import (
	"os"

	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/deployment"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/utils"
)

var config = deployment.Config{
	Name:     "krane-proxy",
	Image:    "biensupernice/proxy",
	Secure:   utils.BoolEnv(constants.EnvProxyDashboardSecure),
	Alias:    []string{os.Getenv(constants.EnvProxyDashboardAlias)},
	Scale:    1,
	Internal: true,
	Env: map[string]string{
		constants.EnvLetsEncryptEmail: os.Getenv(constants.EnvLetsEncryptEmail),
	},
	Volumes: map[string]string{
		"/var/run/docker.sock": "/var/run/docker.sock",
	},
	TargetPort: "8080",
	Ports: map[string]string{
		"80":   "80",
		"443":  "443",
		"8080": "8080",
	},
}

// EnsureNetworkProxy checks that the network proxy has been created and in a running state otherwise will
// attempt to create it. Default behavior is to create the network proxy to allow deployment aliases. This behavior
// can be turned off using the environment variable PROXY_ENABLED which wont create the network proxy if set to false.
func EnsureNetworkProxy() {
	isEnabled := utils.BoolEnv(constants.EnvProxyEnabled)
	if !isEnabled {
		logger.Info("Network proxy not enabled")
		return
	}

	leEmail := os.Getenv(constants.EnvLetsEncryptEmail)
	if config.Secure && leEmail == "" {
		logger.Fatalf("Missing required environment variable %s when running in SECURE mode", constants.EnvLetsEncryptEmail)
	}

	// get containers (if any) for the proxy deployment
	containers, err := deployment.GetContainersByDeployment(config.Name)
	if err != nil {
		logger.Fatalf("Unable to create network proxy, %v", err)
	}

	// create the proxy if no containers are currently up
	if len(containers) == 0 {
		err := createProxy()
		if err != nil {
			// if we cant create the proxy, exit the program
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
				// if we cant create the proxy, exit the program
				logger.Fatalf("Unable to create network proxy, %v", err)
				return
			}
			return
		}
	}

	logger.Debug("Network proxy running")
}

func createProxy() error {
	if err := deployment.Save(config); err != nil {
		return err
	}

	if err := deployment.Run(config.Name); err != nil {
		return err
	}

	logger.Debug("Network proxy deployment started")
	return nil
}
