package longjob

import "sync"

// JobInput The standardized type that all jobs must accept.
type JobInput map[string]interface{}
// JobOutput The standardized output format for all job functions.
type JobOutput interface{}
// JobFunction Standardized job signature for all job functions.
type JobFunction func(payload JobInput) (JobOutput, error)

// Queue Represents a list of inputs received by external services. A go-routine will be spun up for each input
// and the response will be recorded in a map for retrieval by the job creator.
type Queue struct {
	Name           string
	WorkerCount    uint
	WorkerFunction JobFunction
	jobCounter     uint64
	resultQueue    map[uint64]JobResult
	mu             *sync.Mutex
}

// JobResult is a pairing of the outputs of a job function. All jobs must output an object of some kind and/or an error.
type JobResult struct {
	result JobOutput
	err    error
}

// NewQueue Creates a new queue with the specified name. This queue will launch go-routines to run any jobs that
// are added in the queue.
func NewQueue(queueName string, workerCount uint, queueFunction JobFunction) Queue {
	return Queue{
		Name:           queueName,
		WorkerCount:    workerCount,
		WorkerFunction: queueFunction,
		jobCounter:     0,
		resultQueue:    map[uint64]JobResult{},
		mu:             &sync.Mutex{},
	}
}

// Enqueue Adds a job with the specified input to the queue. Returns an instance-specific unique identifier for
// running queries on the job status and results.
func (q *Queue) Enqueue(input JobInput) uint64 {
	q.mu.Lock()
	jobID := q.jobCounter
	q.jobCounter++
	go q.recordJobResult(jobID, q.WorkerFunction, input)
	q.mu.Unlock()
	return jobID
}

// recordJobResult Wrapper that runs the queue JobFunction and records the result in the result map. Function is
// thread safe.
func (q *Queue) recordJobResult(jobID uint64, function JobFunction, input JobInput) {
	output, err := function(input)
	q.mu.Lock()
	q.resultQueue[jobID] = JobResult{
		result: output,
		err:    err,
	}
	q.mu.Unlock()
}

// GetResult Returs the result set from a previous queue job. Assumes that the job is already complete. To test if
// a job is complete, use IsJobDone() first.
func (q *Queue) GetResult(jobID uint64) (JobOutput, error) {
	q.mu.Lock()
	jobResult := q.resultQueue[jobID]
	q.mu.Unlock()
	return jobResult.result, jobResult.err
}

// IsJobDone Checks if results have been reported back for the specified job ID. This function cannot tell the
// difference between a job that has not reported results, and a job that does not exist.
func (q *Queue) IsJobDone(jobID uint64) bool {
	q.mu.Lock()
	_, isPresent := q.resultQueue[jobID]
	q.mu.Unlock()
	return isPresent
}
