package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

// KraneNetworkName : every deployed container will be attached to this network
// TODO: this should be configured somewhere else and passed down when creating docker client / network
const KraneNetworkName = "krane"

// CreateBridgeNetwork : creates docker bridge network with a given name
func (c *DockerClient) CreateBridgeNetwork(ctx *context.Context, name string) (types.NetworkCreateResponse, error) {
	// Check if krane network already exists
	kNet, err := c.GetNetworkByName(ctx, name)
	if err != nil {
		return types.NetworkCreateResponse{}, err
	}
	if kNet.ID != "" {
		return types.NetworkCreateResponse{ID: kNet.ID}, nil
	}

	// If no exisitng network, create it
	options := types.NetworkCreate{
		Driver:         "bridge",
		CheckDuplicate: true,
	}
	return c.NetworkCreate(*ctx, name, options)
}

// GetNetworkByName : find a netwokr by name on this docker host
func (c *DockerClient) GetNetworkByName(ctx *context.Context, name string) (types.NetworkResource, error) {
	// Get all the networks
	options := types.NetworkListOptions{}
	nets, err := c.NetworkList(*ctx, options)
	if err != nil {
		return types.NetworkResource{}, err
	}

	var kNet types.NetworkResource
	for _, net := range nets {
		if net.Name == name {
			kNet = net
			break
		}
	}

	return kNet, nil
}
