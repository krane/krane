package service

import (
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/job"
)

func stopContainerResources(args job.Args) error {
	wf := job.NewWorkflow("StopContainerResources", args)

	wf.With("GetCurrentContainers", getCurrentContainers)
	wf.With("RemoveContainers", stopCurrentContainers)

	return wf.Start()
}

func stopCurrentContainers(args job.Args) error {
	containers := args.GetArg(CurrentContainersJobArgName).(*[]container.KraneContainer)
	for _, c := range *containers {
		err := c.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
