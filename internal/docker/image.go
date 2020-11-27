package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"

	"github.com/biensupernice/krane/internal/constants"
)

type RegistryBasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getRegistryCredentials() string {
	bytes, _ := json.Marshal(RegistryBasicAuth{
		Username: os.Getenv(constants.EnvDockerBasicAuthUsername),
		Password: os.Getenv(constants.EnvDockerBasicAuthPassword),
	})
	return base64.StdEncoding.EncodeToString(bytes)
}

// PullImage : poll docker image from registry
func (c *Client) PullImage(ctx context.Context, registry, image, tag string) (err error) {
	formattedImage := formatImageSourceURL(registry, image, tag)

	options := types.ImagePullOptions{
		All:          true,
		RegistryAuth: getRegistryCredentials(),
	}

	reader, err := c.ImagePull(ctx, formattedImage, options)
	if err != nil {
		return err
	}

	// _, err = io.Copy(nil, reader)
	// if err != nil {
	// 	logrus.Error("Error copying image pull reader to nil dest")
	// }
	err = reader.Close()

	return
}

// RemoveImage : deletes docker image
func (c *Client) RemoveImage(ctx *context.Context, imageID string) ([]types.ImageDelete, error) {
	options := types.ImageRemoveOptions{
		Force:         true, // TODO: was getting race conditions between removing a container and removing the image... couple possible fixes gotta get around to it for now just force remove the images
		PruneChildren: true, // In hopes of keeping the host machine ask light as possible, all child images should be pruned
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
