package job

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddErrorsToJob(t *testing.T) {
	job := Job{Deployment: "test"}

	job.WithError(errors.New("unable to pull image"))
	job.WithError(errors.New("unable to create container"))
	job.WithError(errors.New("unable to Start container"))

	assert.Equal(t, job.Status.Failures[0].Message, "unable to pull image")
	assert.Equal(t, job.Status.Failures[1].Message, "unable to create container")
	assert.Equal(t, job.Status.Failures[2].Message, "unable to Start container")
}
