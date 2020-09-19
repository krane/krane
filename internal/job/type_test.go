package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWhenUnknownJobType(t *testing.T) {
	job := Job{Namespace: "test", Type: ContainerCreate + "x"}
	assert.False(t, isAllowedJobType(job.Type))
}

func TestSuccessWhenKnownJobType(t *testing.T) {
	create := Job{Namespace: "test", Type: ContainerCreate}
	delete := Job{Namespace: "test", Type: ContainerDelete}

	assert.True(t, isAllowedJobType(create.Type))
	assert.True(t, isAllowedJobType(delete.Type))
}
