package service

import (
	"github.com/biensupernice/krane/internal/deployment/container"
	"github.com/biensupernice/krane/internal/job"
)

func deleteContainerResources(args job.Args) error {
	wf := newWorkflow("DeleteContainerResources", args)

	wf.with("GetCurrentContainers", getCurrentContainers)
	wf.with("RemoveContainers", removeCurrContainers)

	return wf.start()
}

func removeCurrContainers(args job.Args) error {
	currContainers := args["currContainers"].(*[]container.Kontainer)
	for _, oldContainer := range *currContainers {
		err := oldContainer.Remove()
		if err != nil {
			return err
		}
	}
	return nil
}
