package deploy

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/biensupernice/krane/result"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DeploySpec struct {
	AppName string           `json:"app" binding:"required"`
	Config  DeploySpecConfig `json:"config" binding:"required"`
}
type DeploySpecConfig struct {
	Repo  string `json:"repo" binding:"required"`
	Image string `json:"image" binding:"required"`
	Tag   string `json:"tag"`
}

func Deploy(spec DeploySpec) (result.Result, error) {
	log.Printf("Deploying %s\n", spec.AppName)

	// Docker Client
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	image := spec.Config.Repo + "/" + spec.Config.Image + ":" + spec.Config.Tag

	log.Printf("Pulling: %s\n", image)

	// Docker pull
	ioreader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, ioreader)
	err = ioreader.Close()
	if err != nil {
		panic(err)
	}

	// Configure Host Port
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8080",
	}

	// Configure Container Port
	containerPort, err := nat.NewPort("tcp", "8080")
	if err != nil {
		panic("Unable to get the port")
	}

	// Bind Host--Container posrts
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	containerConf := &container.Config{Image: image}
	hostConf := &container.HostConfig{PortBindings: portBinding}

	// Docker create container
	resp, err := cli.ContainerCreate(ctx, containerConf, hostConf, nil, "")
	if err != nil {
		panic(err)
	}

	// Docker start container
	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Container %s is started", resp.ID)

	return result.Result{}, nil
}
