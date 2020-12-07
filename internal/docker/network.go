package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// KraneNetworkName : container network for krane containers
const KraneNetworkName = "krane"

// CreateBridgeNetwork : creates docker bridge network
func (c *Client) CreateBridgeNetwork(ctx *context.Context, name string) (types.NetworkCreateResponse, error) {
	n, _ := c.GetNetworkByName(*ctx, name)

	if n.ID != "" {
		return types.NetworkCreateResponse{ID: n.ID}, nil
	}

	return c.NetworkCreate(*ctx, name, types.NetworkCreate{
		Driver:         "bridge",
		CheckDuplicate: true,
	})
}

// GetNetworkByName : find a network by name on this docker host
func (c *Client) GetNetworkByName(ctx context.Context, name string) (types.NetworkResource, error) {
	nets, err := c.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return types.NetworkResource{}, err
	}

	for _, n := range nets {
		if name == n.Name {
			return n, nil
		}
	}

	return types.NetworkResource{}, fmt.Errorf("network %s not found", name)
}

// createNetworkingConfig : create the container network config
func createNetworkingConfig(networkID string) network.NetworkingConfig {
	return network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			KraneNetworkName: {NetworkID: networkID},
		},
	}
}
