package service

import (
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/job"
)

func deleteContainerResources(args job.Args) error {
	namespace := args["namespace"]
	logrus.Printf("Deleting deployment workflow for %s", namespace)
	return nil
}
