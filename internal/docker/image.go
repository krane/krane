package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
)

// PullImage pulls a container image from a registry onto the host machine
func (c *Client) PullImage(image string, tag string, registry RegistryCredentials) (io.Reader, error) {
	ctx := context.Background()
	defer ctx.Done()

	ref := createImageRef(registry.URL, image, tag)
	return c.ImagePull(ctx, ref, types.ImagePullOptions{
		All:          false,
		RegistryAuth: Base64RegistryCredentials(registry.Username, registry.Password),
	})
}

// RemoveImage removes a docker image from the host machine
func (c *Client) RemoveImage(ctx *context.Context, imageID string) ([]types.ImageDelete, error) {
	options := types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}
	return c.ImageRemove(*ctx, imageID, options)
}

// createImageRef returns a formatted docker image url
func createImageRef(registry, image, tag string) string {
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s/%s:%s", registry, image, tag)
}
