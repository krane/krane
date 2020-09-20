package job

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartJob(t *testing.T) {
	j := Job{}
	j.start()
	assert.Equal(t, Started, j.State)
	assert.True(t, time.Now().Unix() >= j.StartTime)
}

func TestEndJob(t *testing.T) {
	j := Job{}

	j.start()
	assert.Equal(t, Started, j.State)
	assert.True(t, time.Now().Unix() >= j.StartTime)

	j.end()
	assert.Equal(t, Complete, j.State)
	assert.True(t, time.Now().Unix() >= j.EndTime)
}

func TestJobNotStartedOnCallToJobEnd(t *testing.T) {
	j := Job{}

	j.end()
	assert.NotEqual(t, Complete, j.State)

	j.start()
	assert.Equal(t, Started, j.State)
	assert.True(t, time.Now().Unix() >= j.StartTime)

	j.end()
	assert.Equal(t, Complete, j.State)
	assert.True(t, time.Now().Unix() >= j.EndTime)
}
