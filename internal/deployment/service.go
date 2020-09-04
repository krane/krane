package deployment

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"

	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/spec"
)

type Deployment struct {
	Alias      string            `json:"alias"`
	Spec       spec.Spec         `json:"spec"`
	Containers []types.Container `json:"containers"`
}

// Start :  a deployment creating any docker resources required
func (d Deployment) Start(props map[string]string) error {

	// Delete containers first otherwise conflicting ports will occur. In the future it would be n
	// eat to have this do some kind of rolling deployment. Messing with the ports and traefik labels might make this possible
	err := d.deleteDockerResources()
	if err != nil {
		return err
	}

	// If the incoming tag is different from the one on the Spec, update the spec
	tag := props["tag"]
	if tag != d.Spec.Config.Tag {
		d.Spec.Config.Tag = tag
		err := d.Spec.UpdateSpec(d.Spec.Name)
		if err != nil {
			return err
		}
	}

	err = d.startDockerResources()
	if err != nil {
		return err
	}
	return nil
}

// Start :  a deployment. Usually consists of stopping all the containers.
func (d Deployment) Stop(props map[string]string) error {
	err := d.stopDockerResources()
	if err != nil {
		return err
	}

	return nil
}

// Delete : a deployment. usually deleting all the containers and removes the deployment spec, and alias.
func (d Deployment) Delete(props map[string]string) error {
	// Delete docker resources
	err := d.deleteDockerResources()
	if err != nil {
		return err
	}

	// Get any active containers
	containers, err := GetDeploymentContainers(d.Spec.Name)
	if err != nil {
		return err
	}

	// If any containers are still created, do not delete the spec or alias..
	if len(containers) != 0 {
		return errors.New("Cannot continue with delete. Container count for this deployment is greater than 0")
	}

	// Delete Spec
	d.Spec.Delete()

	// Delete Alias
	d.DeleteAlias(map[string]string{"alias": d.Alias})

	return nil
}

func GetDeployment(deploymentId string) (Deployment, error) {
	// Find spec
	spec, err := spec.GetOne(deploymentId)
	if err != nil {
		return Deployment{}, err
	}

	ctx := context.Background()

	// Find containers
	containers, err := docker.GetContainers(&ctx, deploymentId)
	if err != nil {
		return Deployment{}, err
	}

	ctx.Done()

	// Find alias
	alias, err := GetDeploymentAlias(deploymentId)
	if err != nil {
		return Deployment{}, err
	}

	// Build deployment struct
	d := *&Deployment{
		Alias:      alias,
		Spec:       spec,
		Containers: containers,
	}

	return d, nil
}

func GetDeploymentContainers(deploymentId string) ([]types.Container, error) {
	ctx := context.Background()
	defer ctx.Done()

	// Find containers
	return docker.GetContainers(&ctx, deploymentId)
}

func GetDeployments() ([]Deployment, error) {
	// Get all specs
	specs, err := spec.GetAll()
	if err != nil {
		return make([]Deployment, 0), err
	}

	// Get all deployments using the spec name
	deployments := make([]Deployment, 0)
	for _, s := range specs {
		d, err := GetDeployment(s.Name)
		if err != nil {
			return deployments, err
		}

		deployments = append(deployments, d)
	}

	return deployments, nil
}
