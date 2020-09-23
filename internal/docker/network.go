package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// KraneNetworkName : every deployed container will be attached to this network
const KraneNetworkName = "krane"

func makeNetworkingConfig(networkID string) network.NetworkingConfig {
	endpointConf := map[string]*network.EndpointSettings{"krane": &network.EndpointSettings{NetworkID: networkID}}
	return network.NetworkingConfig{
		EndpointsConfig: endpointConf,
	}
}

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
