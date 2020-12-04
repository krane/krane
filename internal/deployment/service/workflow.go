package service

import (
	"github.com/biensupernice/krane/internal/job"
	"github.com/biensupernice/krane/internal/logger"
)

type workflow struct {
	name string
	args job.Args
	head *step
	curr *step
}

type step struct {
	name string
	fn   stepFn
	next *step
}

type stepFn func(arg job.Args) error

func newWorkflow(name string, args job.Args) workflow {
	return workflow{name: name, args: args}
}

// with : is a method on  workflow for adding new steps.
// The workflow is a linked list of execution steps, stepFn.
func (wf *workflow) with(name string, handler stepFn) {
	s := &step{name: name, fn: handler}
	if wf.head == nil {
		wf.head = s
	} else {
		currStep := wf.head
		for currStep.next != nil {
			currStep = currStep.next
		}
		currStep.next = s
	}
}

// start : executing every step in a workflow.
// It returns an error if any step in the workflow errors out.
// The idea behind a workflow is that you want to execute multiple
// steps but want things to maintain readability and flow.
// So every step is just a function that should received job.Args and return an error.
func (wf *workflow) start() error {
	wf.curr = wf.head

	// run every step starting from the head of the workflow
	// until there aren't any steps left to run
	for wf.curr != nil {
		logger.Debugf("Running Workflow %s | Step %s", wf.name, wf.curr.name)

		// execute every step passing down args
		err := wf.curr.fn(wf.args)
		if err != nil {
			// if any step fails, the workflow
			// stops executing further steps
			return err
		}

		wf.curr = wf.next()
	}

	return nil
}

func (wf *workflow) next() *step {
	if wf.curr == nil {
		return nil
	}

	wf.curr = wf.curr.next
	return wf.curr
}
