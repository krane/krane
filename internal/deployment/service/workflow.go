package service

import "github.com/docker/docker/api/types"

type workflow struct {
	name string
	head *step
	curr *step
}

type step struct {
	name string
	fn   stepFn
	next *step
}

type stepFn func(arg types.Arg) (error)

func newWorkflow(name string) workflow { return workflow{name: name} }

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

func (wf *workflow) start() *step {
	wf.curr = wf.head
	return wf.curr
}

func (wf *workflow) next() *step {
	wf.curr = wf.curr.next
	return wf.curr
}
