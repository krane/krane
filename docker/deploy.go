package docker

import (
	"context"
	"log"
)

// DeploySpec : spec to deploy and app
type DeploySpec struct {
	Name   string           `json:"name" binding:"required"`
	Config DeploySpecConfig `json:"config" binding:"required"`
}

// DeploySpecConfig : config for deploying an app
type DeploySpecConfig struct {
	Registry      string `json:"repo"`                     // Docker registry url
	Image         string `json:"image" binding:"required"` // DOcker image name
	Tag           string `json:"tag"`                      // Docker image tag
	HostPort      string `json:"host_port"`                // Port to bind to  host machine from the container
	ContainerPort string `json:"container_port"`           // POrt to expose from the container
}

// Deploy : docker container
func Deploy(spec DeploySpec) (containerID string, err error) {

	// Set docker registry if not provided
	if spec.Config.Registry == "" {
		spec.Config.Registry = "registry.hub.docker.com"
	}

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

	log.Printf("Deploying %s\n", spec.Name)

	// Get docker client
	_, err = New()
	if err != nil {
		return
	}

	// Start deployment context
	ctx := context.Background()

	// Format docker image url source
	img := FormatImageSourceURL(spec.Config.Registry, spec.Config.Image, spec.Config.Tag)
	log.Printf("Pulling image: %s\n", img)

	// Pull docker image
	err = PullImage(&ctx, img)
	if err != nil {
		return
	}

	// Create docker container
	createContainerResp, err := CreateContainer(&ctx, img, spec.Name, spec.Config.HostPort, spec.Config.ContainerPort)
	containerID = createContainerResp.ID
	if err != nil {
		return
	}

	// Docker start container
	err = StartContainer(&ctx, containerID)
	if err != nil {
		// If error starting the container, remove it
		RemoveContainer(&ctx, containerID)
		return
	}

	log.Printf("Deployed %s - 📦 %s\n", spec.Name, containerID)

	return
}
