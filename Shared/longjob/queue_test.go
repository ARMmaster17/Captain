package longjob

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreatesANamedQueue(t *testing.T) {
	queue := NewQueue("shared-queue", 1)
	assert.Equal(t, "shared-queue", queue.Name)
	assert.Equal(t, uint(1), queue.WorkerCount)
}
