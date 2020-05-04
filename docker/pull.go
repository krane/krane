package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Pull docker image
func Pull(ctx context.Context, cli client.Client, image string) error {
	ioreader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, ioreader)
	err = ioreader.Close()
	if err != nil {
		return err
	}

	return nil
}
