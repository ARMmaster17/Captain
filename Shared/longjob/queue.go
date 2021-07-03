package longjob

import "sync"

type JobInput map[string]interface{}
type JobOutput interface{}
type JobFunction func(payload JobInput) (JobOutput, error)

type Queue struct {
	Name string
	WorkerCount uint
	WorkerFunction JobFunction
	jobCounter uint64
	resultQueue map[uint64]JobResult
	mu *sync.Mutex
}

type JobResult struct {
	result JobOutput
	err error
}

// NewQueue Creates a new queue with the specified name. This queue will launch go-routines to run any jobs that
// are added in the queue.
func NewQueue(queueName string, workerCount uint, queueFunction JobFunction) Queue {
	return Queue{
		Name: queueName,
		WorkerCount: workerCount,
		WorkerFunction: queueFunction,
		jobCounter: 0,
		resultQueue: map[uint64]JobResult{},
		mu: &sync.Mutex{},
	}
}

// Enqueue Adds a job with the specified input to the queue. Returns an instance-specific unique identifier for
// running queries on the job status and results.
func (q *Queue) Enqueue(input JobInput) uint64 {
	q.mu.Lock()
	jobId := q.jobCounter
	q.jobCounter++
	go q.recordJobResult(jobId, q.WorkerFunction, input)
	q.mu.Unlock()
	return jobId
}

// recordJobResult Wrapper that runs the queue JobFunction and records the result in the result map. Function is
// thread safe.
func (q *Queue) recordJobResult(jobId uint64, function JobFunction, input JobInput) {
	output, err := function(input)
	q.mu.Lock()
	q.resultQueue[jobId] = JobResult{
		result: output,
		err:    err,
	}
	q.mu.Unlock()
}

// GetResult Returs the result set from a previous queue job. Assumes that the job is already complete. To test if
// a job is complete, use IsJobDone() first.
func (q *Queue) GetResult(jobId uint64) (JobOutput, error) {
	q.mu.Lock()
	jobResult := q.resultQueue[jobId]
	q.mu.Unlock()
	return jobResult.result, jobResult.err
}

// IsJobDone Checks if results have been reported back for the specified job ID. This function cannot tell the
// difference between a job that has not reported results, and a job that does not exist.
func (q *Queue) IsJobDone(jobId uint64) bool {
	q.mu.Lock()
	_, isPresent := q.resultQueue[jobId]
	q.mu.Unlock()
	return isPresent
}
