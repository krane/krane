package job

import (
	"github.com/biensupernice/krane/internal/logger"
)

type Workflow struct {
	name string
	args interface{}
	head *Step
	curr *Step
}

type Step struct {
	name string
	fn   GenericHandler
	next *Step
}

// NewWorkflow : creates a new workflow
func NewWorkflow(name string, args interface{}) Workflow {
	return Workflow{name: name, args: args}
}

// With : add new step to a workflow
func (wf *Workflow) With(name string, handler GenericHandler) {
	s := &Step{name: name, fn: handler}
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

// Start : executes every Step in a Workflow
// returns an error if any Step in the Workflow errors out.
func (wf *Workflow) Start() error {
	wf.curr = wf.head

	// run every Step starting from the head of the Workflow
	// until there aren't any steps left to run
	for wf.curr != nil {
		logger.Debugf("Running Workflow %s | Step %s", wf.name, wf.curr.name)

		// execute every Step passing down args
		err := wf.curr.fn(wf.args)
		if err != nil {
			// if any Step fails, the Workflow
			// stops executing further steps
			return err
		}

		wf.curr = wf.next()
	}

	return nil
}

// next : execute the next Step (if any) in the Workflow
func (wf *Workflow) next() *Step {
	if wf.curr == nil {
		return nil
	}

	wf.curr = wf.curr.next
	return wf.curr
}
