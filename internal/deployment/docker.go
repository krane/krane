package deployment

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/docker"
)

func (d *Deployment) startDockerResources() error {

	ctx := context.Background()

	client := docker.GetClient()

	// Create formatted url to fetch docker image - http://registry/image:tag
	image := docker.FormatImageSourceURL(d.Spec.Config.Registry, d.Spec.Config.Image, d.Spec.Config.Tag)

	logrus.Debugf("[%s] -> Pulling %s", d.Spec.Name, image)
	err := client.PullImage(&ctx, image)
	if err != nil {
		logrus.Errorf("[%s] -> Error pulling %s: %s", d.Spec.Name, image, err.Error())
		return err
	}

	// Container Labels added for identifying krane managed containers
	routingLabel := fmt.Sprintf("traefik.status.routers.%s.rule", d.Spec.Name)
	routingValue := fmt.Sprintf("Host(`%s`)", d.Alias)

	ruleLabel := fmt.Sprintf("traefik.%s.rule", d.Spec.Name)
	ruleKey := fmt.Sprintf("Host:%s", d.Alias)

	labels := map[string]string{
		docker.KraneContainerLabelName: d.Spec.Name,
		"traefik.enable":               "true",
		routingLabel:                   routingValue,
		ruleLabel:                      ruleKey,
	}

	if d.Spec.Config.ContainerPort != "" {
		loadBalancerPortLabel := fmt.Sprintf("traefik.status.services.%s.loadbalancer.server.port", d.Spec.Name) // If the container exposes more than 2 ports this value should be set using the container_port property in krane.json to tell traefik to use a specific port for routing
		labels[loadBalancerPortLabel] = d.Spec.Config.ContainerPort
	}

	// Krane Network ID to connect the container
	net, err := client.GetNetworkByName(&ctx, "krane")
	if err != nil {
		logrus.Errorf("[%s] -> Error get network id: %s", d.Spec.Name, err.Error())
		return err
	}

	// Create the config used to create the container
	containerConf := &docker.CreateContainerConfig{
		Name:          d.Spec.Name,
		Image:         image,
		NetworkID:     net.ID,
		HostPort:      d.Spec.Config.HostPort,
		ContainerPort: d.Spec.Config.ContainerPort,
		Env:           d.Spec.Config.Env,
		Labels:        labels,
		Volumes:       d.Spec.Config.Volumes,
	}

	// Create the container on the docker host machine
	logrus.Debugf("[%s] -> Creating container %s", d.Spec.Name, d.Spec.Name)
	createContainerResp, err := client.CreateContainer(&ctx, containerConf)
	if err != nil {
		return err
	}

	containerID := createContainerResp.ID
	logrus.Debugf("[%s] -> Container created with id %s", d.Spec.Name, containerID)

	// Connect container to network
	err = client.ConnectContainerToNetwork(&ctx, net.ID, containerID)
	if err != nil {
		logrus.Errorf("[%s] -> Error unable to connect container to the docker network: %s", d.Spec.Name, err.Error())
		return err
	}

	// Start the already created containers
	logrus.Debugf("[%s] -> Starting container %s", d.Spec.Name, containerID)
	err = client.StartContainer(&ctx, containerID)
	if err != nil {
		logrus.Errorf("[%s] -> Error unable to start container %s: %s", d.Spec.Name, containerID, err.Error())
		err := client.RemoveContainer(&ctx, containerID)
		if err != nil {
			logrus.Errorf("[%s] -> Error unable to remove container %s: %s", d.Spec.Name, containerID, err.Error())
			return err
		}
		return err
	}

	logrus.Debugf("[%s] -> Container %s started", d.Spec.Name, containerID)

	ctx.Done()

	return nil
}

func (d *Deployment) stopDockerResources() error {

	ctx := context.Background()

	client := docker.GetClient()
	for _, container := range d.Containers {
		err := client.StopContainer(&ctx, container.ID)
		if err != nil {
			logrus.Errorf("[%s] -> Error stopping containers %s", d.Spec.Name, err.Error())
			return err
		}
	}

	ctx.Done()

	logrus.Debugf("[%s] -> Stopped containers", d.Spec.Name)
	return nil
}

func (d *Deployment) deleteDockerResources() error {
	ctx := context.Background()

	client := docker.GetClient()
	for _, container := range d.Containers {
		// Stop the container
		err := client.StopContainer(&ctx, container.ID)
		if err != nil {
			logrus.Errorf("[%s] -> Unable to stop %s - %s", d.Spec.Name, container.ID, err.Error())
			return err
		}
		logrus.Debugf("[%s] -> Stopped container %s", d.Spec.Name, container.ID)

		// Remove the container
		err = client.RemoveContainer(&ctx, container.ID)
		if err != nil {
			logrus.Errorf("[%s] -> Unable to remove %s - %s", d.Spec.Name, container.ID, err.Error())
			return err
		}
		logrus.Debugf("[%s] -> Removed container %s", d.Spec.Name, container.ID)
	}

	// Remove the image(s)
	for _, container := range d.Containers {
		_, err := client.RemoveImage(&ctx, container.ImageID)
		if err != nil {
			logrus.Errorf("[%s] -> Unable to remove image %s - %s", d.Spec.Name, container.ImageID, err.Error())
			return err
		}
		logrus.Debugf("[%s] -> Removed image %s", d.Spec.Name, container.ImageID)
	}

	ctx.Done()

	return nil
}
