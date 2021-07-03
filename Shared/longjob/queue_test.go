package longjob

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_CreatesANamedQueue(t *testing.T) {
	queue := NewQueue("shared-queue", 1, JobFunction(nil))
	assert.Equal(t, "shared-queue", queue.Name)
	assert.Equal(t, uint(1), queue.WorkerCount)
}

func Test_AddJobToQueue(t *testing.T) {
	queue := NewQueue("shared-queue", 1, func(payload JobInput) (JobOutput, error) {
		time.Sleep(1 * time.Second)
		return "", nil
	})
	assert.Equal(t, "shared-queue", queue.Name)
	assert.Equal(t, uint(1), queue.WorkerCount)
	var sampleInput map[string]interface{}
	queue.Enqueue(sampleInput)
}

func Test_GetsJobIdAfterEnqueue(t *testing.T) {
	queue := NewQueue("shared-queue", 1, func(payload JobInput) (JobOutput, error) {
		time.Sleep(1 * time.Second)
		return "", nil
	})
	assert.Equal(t, "shared-queue", queue.Name)
	assert.Equal(t, uint(1), queue.WorkerCount)
	var sampleInput map[string]interface{}
	result := queue.Enqueue(sampleInput)
	assert.Equal(t, uint64(0), result)
}

func Test_GetsIncrementingJobIdAfterEnqueue(t *testing.T) {
	queue := NewQueue("shared-queue", 1, func(payload JobInput) (JobOutput, error) {
		time.Sleep(1 * time.Second)
		return "", nil
	})
	assert.Equal(t, "shared-queue", queue.Name)
	assert.Equal(t, uint(1), queue.WorkerCount)
	var sampleInput map[string]interface{}
	result := queue.Enqueue(sampleInput)
	assert.Equal(t, uint64(0), result)
	result2 := queue.Enqueue(sampleInput)
	assert.Equal(t, uint64(1), result2)
}

func Test_RecordsJobResults(t *testing.T) {
	queue := NewQueue("shared-queue", 1, func(payload JobInput) (JobOutput, error) {
		return "test", nil
	})
	assert.Equal(t, "shared-queue", queue.Name)
	assert.Equal(t, uint(1), queue.WorkerCount)
	var sampleInput map[string]interface{}
	jobId := queue.Enqueue(sampleInput)
	assert.Equal(t, uint64(0), jobId)
	for {
		if queue.IsJobDone(jobId) {
			break
		}
	}
	jobResult, err := queue.GetResult(jobId)
	assert.Equal(t, "test", jobResult)
	assert.NoError(t, err)
}

func Test_ReportsIfJobIsDone(t *testing.T) {
	queue := NewQueue("shared-queue", 1, func(payload JobInput) (JobOutput, error) {
		time.Sleep(1 * time.Second)
		return "test", nil
	})
	assert.Equal(t, "shared-queue", queue.Name)
	assert.Equal(t, uint(1), queue.WorkerCount)
	var sampleInput map[string]interface{}
	jobId := queue.Enqueue(sampleInput)
	assert.Equal(t, uint64(0), jobId)
	assert.Equal(t, false, queue.IsJobDone(jobId))
}
