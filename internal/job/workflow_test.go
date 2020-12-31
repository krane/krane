package job

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflowWithNoStepsDoesntError(t *testing.T) {
	wf := NewWorkflow("noSteps", nil)
	err := wf.Start()
	assert.Nil(t, err)
}

func TestWorkflowWithNoStepsOnCallNextReturnsNil(t *testing.T) {
	wf := NewWorkflow("noSteps", nil)
	assert.Nil(t, wf.next())
}

func TestWorkflowWithNSteps(t *testing.T) {
	// This test is formatted a bit weird to test
	// that the steps in a Workflow are executed.
	// A pointer to a variable stepCount is passed
	// into the Step arguments and every Step increments
	// x by 1. To test x gets executed we do some pointer
	// manipulation for testing purposes. Its not really a great
	// idea in practice to be passing pointers to values since
	// erroneous steps could change the value and error out entire workflows.

	// the variable under tests we wanna increase
	x := 0

	// Step function used to increment x
	incX := func(args interface{}) error {
		x := args.(map[string]*int)["stepCount"]
		*x++
		return nil
	}

	// argument passed to every Step in the Workflow
	args := map[string]*int{"stepCount": &x}

	wf := NewWorkflow("testSteps", args)

	stepCount := 20
	for i := 0; i < stepCount; i++ {
		stepName := fmt.Sprintf("step_%d", i)

		// add new Step to the Workflow
		// in this example we just want to create stepCount amount
		// of steps and incremet x, stepCount amount of times.
		wf.With(stepName, incX)
	}

	// Start the Workflow
	err := wf.Start()

	assert.Nil(t, err)
	assert.Equal(t, stepCount, *args["stepCount"])
}

func TestWorkflowError(t *testing.T) {
	wf := NewWorkflow("testWorkflowError", nil)

	step := func(args interface{}) error {
		if args == nil {
			return errors.New("Step args cannot be nil")
		}
		return nil
	}

	wf.With("VerifyArgsNotNil", step)

	err := wf.Start()

	assert.Error(t, err)
	assert.Equal(t, "Step args cannot be nil", err.Error())
}
