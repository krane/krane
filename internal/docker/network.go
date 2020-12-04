package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// KraneNetworkName : krane managed network
const KraneNetworkName = "krane"

func makeNetworkingConfig(networkID string) network.NetworkingConfig {
	endpointConf := map[string]*network.EndpointSettings{"krane": {NetworkID: networkID}}
	return network.NetworkingConfig{
		EndpointsConfig: endpointConf,
	}
}

// CreateBridgeNetwork : creates docker bridge network with a given name
func (c *Client) CreateBridgeNetwork(ctx *context.Context, name string) (types.NetworkCreateResponse, error) {
	// Check if krane network already exists
	kNet, err := c.GetNetworkByName(*ctx, name)
	if err != nil {
		return types.NetworkCreateResponse{}, err
	}
	if kNet.ID != "" {
		return types.NetworkCreateResponse{ID: kNet.ID}, nil
	}

	// If no existing network, create it
	options := types.NetworkCreate{
		Driver:         "bridge",
		CheckDuplicate: true,
	}
	return c.NetworkCreate(*ctx, name, options)
}

// GetNetworkByName : find a network by name on this docker host
func (c *Client) GetNetworkByName(ctx context.Context, name string) (types.NetworkResource, error) {
	// Get all the networks
	options := types.NetworkListOptions{}
	nets, err := c.NetworkList(ctx, options)
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
