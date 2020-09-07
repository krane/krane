package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

// PullImage : poll docker image from registry
func (c *DockerClient) PullImage(ctx *context.Context, image string) (err error) {
	options := types.ImagePullOptions{
		RegistryAuth: "", // RegistryAuth is the base64 encoded credentials for the registry
	}

	reader, err := c.ImagePull(*ctx, image, options)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, reader)
	err = reader.Close()

	return
}

// RemoveImage : deletes docker image
func (c *DockerClient) RemoveImage(ctx *context.Context, imageID string) ([]types.ImageDelete, error) {
	options := types.ImageRemoveOptions{
		Force:         true, // TODO: was getting race conditions between removing a container and removing the image... couple possible fixes gotta get around to it for now just force remove the images
		PruneChildren: true, // In hopes of keeping the host machine ask light as possible, all child images should be pruned
	}
	return c.ImageRemove(*ctx, imageID, options)
}

// FormatImageSourceURL : format into appropriate docker image url
func FormatImageSourceURL(
	repo string,
	imageName string,
	tag string) string {
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s/%s:%s", repo, imageName, tag)
}
