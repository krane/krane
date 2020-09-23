package handler

import (
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/job"
)

type RunDeploymentArgs struct {
	Config config.Config
}

func RunDeploymentHandler(args job.Args) error {
	return nil
}
