package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// KraneNetworkName is the network used for Krane containers
const KraneNetworkName = "krane"

// CreateBridgeNetwork creates a docker bridge network
func (c *Client) CreateBridgeNetwork(ctx *context.Context, name string) (types.NetworkCreateResponse, error) {
	n, _ := c.GetNetworkByName(name)

	if n.ID != "" {
		return types.NetworkCreateResponse{ID: n.ID}, nil
	}

	return c.NetworkCreate(*ctx, name, types.NetworkCreate{
		Driver:         "bridge",
		CheckDuplicate: true,
	})
}

// GetNetworkByName returns the network (if it exist) from the docker host
func (c *Client) GetNetworkByName(name string) (types.NetworkResource, error) {
	ctx := context.Background()
	defer ctx.Done()

	networks, err := c.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return types.NetworkResource{}, err
	}

	for _, net := range networks {
		if name == net.Name {
			return net, nil
		}
	}

	return types.NetworkResource{}, fmt.Errorf("network %s not found", name)
}

// createNetworkingConfig create the container network config
func createNetworkingConfig(networkID string) network.NetworkingConfig {
	return network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			KraneNetworkName: {
				NetworkID: networkID,
			},
		},
	}
}
