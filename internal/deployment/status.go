package deployment

import (
	"context"
	"io"
	"os"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/logger"
)

const (
	// ReadyStatus : Deployment is ready
	ReadyStatus = "Ready"
	// DeployingStatus : Deployment is in progress
	DeployingStatus = "Deploying"
	// FailedStatus : Deployment has failed
	FailedStatus = "Failed"
)

// GetStatus : of a deployment by name returns -- "Ready, Deploying, Failed"
func GetStatus(ctx *context.Context, name string) (status string, err error) {
	logger.Debugf("Getting status for deployment - %s", name)
	containerID := "b0d44cf4956c"
	stats, err := docker.GetContainerStatus(ctx, containerID, false)
	if err != nil {
		return
	}

	io.Copy(os.Stdout, stats.Body)
	err = stats.Body.Close()

	return
}
