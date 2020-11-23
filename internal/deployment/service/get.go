package service

import (
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/job"
)

func getCurrentContainers(args job.Args) error {
	cfg := args["kconfig"].(kconfig.Kconfig)
	containers, err := container.GetKontainersByNamespace(cfg.Name)
	if err != nil {
		return err
	}
	args["currContainers"] = &containers
	return nil
}
