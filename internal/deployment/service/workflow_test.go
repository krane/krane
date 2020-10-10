package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biensupernice/krane/internal/job"
)

func TestWorkflowWithNoStepsDoesntError(t *testing.T) {
	wf := newWorkflow("noSteps", nil)
	err := wf.start()
	assert.Nil(t, err)
}

func TestWorkflowWithNoStepsOnCallNextReturnsNil(t *testing.T) {
	wf := newWorkflow("noSteps", nil)
	assert.Nil(t, wf.next())
}

func TestWorkflowWithNSteps(t *testing.T) {
	// This test is formatted a bit weird to test
	// that the steps in a workflow are executed.
	// A pointer to a variable stepCount is passed
	// into the step arguments and every step increments
	// x by 1. To test x gets executed we do some pointer
	// manipulation for testing purposes. Its not really a great
	// idea in practice to be passing pointers to values since
	// erroneous steps could change the value and error out entire workflows.

	// the variable under tests we wanna increase
	x := 0

	// step function used to increment x
	incX := func(args job.Args) error {
		x := args["stepCount"].(*int)
		*x++
		return nil
	}

	// argument passed to every step in the workflow
	args := job.Args{"stepCount": &x}

	wf := newWorkflow("testSteps", args)

	stepCount := 20
	for i := 0; i < stepCount; i++ {
		stepName := fmt.Sprintf("step_%d", i)

		// add new step to the workflow
		// in this example we just want to create stepCount amount
		// of steps and incremet x, stepCount amount of times.
		wf.with(stepName, incX)
	}

	// start the workflow
	err := wf.start()

	assert.Nil(t, err)
	assert.Equal(t, stepCount, *args["stepCount"].(*int))
}

func TestWorkflowError(t *testing.T) {
	wf := newWorkflow("testWorkflowError", nil)

	step := func(args job.Args) error {
		if args == nil {
			return errors.New("step args cannot be nil")
		}
		return nil
	}

	wf.with("VerifyArgsNotNil", step)

	err := wf.start()

	assert.Error(t, err)
	assert.Equal(t, "step args cannot be nil", err.Error())
}
