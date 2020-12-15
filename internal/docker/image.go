package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

// PullImage : pull docker image from registry
func (c *Client) PullImage(ctx context.Context, registry, image, tag string) (err error) {
	formattedImage := formatImageSourceURL(registry, image, tag)

	options := types.ImagePullOptions{
		All:          false,
		RegistryAuth: Base64DockerRegistryCredentials(),
	}

	reader, err := c.ImagePull(ctx, formattedImage, options)
	if err != nil {
		return err
	}

	// TODO: dont output to standard out
	// send as a deployment event
	io.Copy(os.Stdout, reader)

	err = reader.Close()

	return
}

// RemoveImage : remove docker image
func (c *Client) RemoveImage(ctx *context.Context, imageID string) ([]types.ImageDelete, error) {
	options := types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}
	return c.ImageRemove(*ctx, imageID, options)
}

// formatImageSourceURL : format into appropriate docker image url
func formatImageSourceURL(registry, image, tag string) string {
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s/%s:%s", registry, image, tag)
}
