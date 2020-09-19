package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueCapacityOnNewQueueInstance(t *testing.T) {
	qSize := 10
	queue := NewJobQueue(uint(qSize))
	assert.Equal(t, qSize, cap(queue))
	assert.Equal(t, 0, len(queue))
}

func TestReturnSameQueueInstance(t *testing.T) {
	qSize := 10
	q1 := NewJobQueue(uint(qSize))
	q2 := NewJobQueue(uint(qSize))
	assert.Exactly(t, &q1, &q2)
}
