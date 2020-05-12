package docker

import (
	"context"
	"log"
)

// DeploySpec : spec to deploy and app
type DeploySpec struct {
	AppName string           `json:"app" binding:"required"`
	Config  DeploySpecConfig `json:"config" binding:"required"`
}

// DeploySpecConfig : config for deploying an app
type DeploySpecConfig struct {
	Repo          string `json:"repo" binding:"required"`
	Image         string `json:"image" binding:"required"`
	Tag           string `json:"tag"`
	HostPort      string `json:"host_port"`
	ContainerPort string `json:"container_port"`
}

// Deploy docker container
func Deploy(spec DeploySpec) error {

	// Set image tag to `latest` if not provided
	if spec.Config.Tag == "" {
		spec.Config.Tag = "latest"
	}

	// Set container host port to `8080` if not provided
	if spec.Config.HostPort == "" {
		spec.Config.HostPort = "8080"
	}

	// Set container port to `8080` if not provided
	if spec.Config.ContainerPort == "" {
		spec.Config.ContainerPort = "8080"
	}

	log.Printf("Deploying %s\n", spec.AppName)

	// Create docker client
	_, err := New()
	if err != nil {
		log.Printf("Unable to create docker client %s\n", err.Error())
		return err
	}

	img := FormatImageSourceURL(spec.Config.Repo, spec.Config.Image, spec.Config.Tag)

	ctx := context.Background() // deployment context

	log.Printf("Pulling image: %s\n", img)

	// Pull docker image
	err = PullImage(&ctx, img)
	if err != nil {
		log.Printf("Unable to pull image %s - %s\n", img, err.Error())
		return err
	}

	// Create docker container
	createContainerResp, err := CreateContainer(&ctx, img, "", spec.Config.HostPort, spec.Config.ContainerPort)
	containerID := createContainerResp.ID
	if err != nil {
		log.Printf("Unable to create container for image %s - %s\n", img, err.Error())
		return nil
	}

	// Docker start container
	err = StartContainer(&ctx, containerID)
	if err != nil {
		log.Printf("Unable to start container %s - %s", containerID, err.Error())
		return nil
	}
	log.Printf("Deployed %s - 📦 %s\n", spec.AppName, containerID)

	return nil
}
