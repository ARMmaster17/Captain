package main

// Lots of this code was initially borrowed from http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/.

import (
	"fmt"
	"os"
)

type BuilderDispatcher struct {
	BuilderPool chan chan BuilderJob
}

var GlobalBuilderDispatcher BuilderDispatcher

func InitBuilderDispatcher(maxBuilders int) {
	pool := make(chan chan BuilderJob, maxBuilders)
	GlobalBuilderDispatcher = BuilderDispatcher{BuilderPool: pool}
	for i := 0; i < maxBuilders; i++ {
		builder := NewBuilder(GlobalBuilderDispatcher.BuilderPool)
		builder.Start()
	}
	go GlobalBuilderDispatcher.Dispatch()
}

func (d *BuilderDispatcher) Dispatch() {
	for {
		select {
			case job := <-BuilderJobQueue:
				go func(job BuilderJob) {
					jobChannel := <-d.BuilderPool
					jobChannel <- job
				} (job)
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////
var (
	MaxWorker = os.Getenv("MAX_WORKERS")
	MaxQueue  = os.Getenv("MAX_QUEUE")
)

// Job represents the job to be run
type BuilderJob struct {
	Payload BuilderPayload
}

type BuilderPayload struct {
	NewVMID int
	Plane Plane
}

// A buffered channel that we can send work requests on.
var BuilderJobQueue chan BuilderJob

// Worker represents the worker that executes the job
type Builder struct {
	WorkerPool  chan chan BuilderJob
	JobChannel  chan BuilderJob
	quit    	chan bool
}

func NewBuilder(workerPool chan chan BuilderJob) Builder {
	return Builder{
		WorkerPool: workerPool,
		JobChannel: make(chan BuilderJob),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Builder) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				job.Payload.Plane.Validate()
				newPlane := Plane{
					Num: f.getNextNum(i),
					FormationID: int(f.ID),
				}
				result := db.Save(&newPlane)
				if result.Error != nil {
					return fmt.Errorf("unable to update formation with new planes with error: %w", result.Error)
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Builder) Stop() {
	go func() {
		w.quit <- true
	}()
}