package main

// Lots of this code was initially borrowed from http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/.

import (
	"fmt"
	"github.com/ARMmaster17/Captain/db"
	"github.com/rs/zerolog/log"
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
		builder.Start(i)
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

// Job represents the job to be run
type BuilderJob struct {
	Payload BuilderPayload
}

type BuilderPayload struct {
	Plane Plane
}

// A buffered channel that we can send work requests on.
var BuilderJobQueue chan BuilderJob

// Worker represents the worker that executes the job
type Builder struct {
	WorkerPool  chan chan BuilderJob
	JobChannel  chan BuilderJob
	quit    	chan bool
	ID			int
}

func NewBuilder(workerPool chan chan BuilderJob) Builder {
	return Builder{
		WorkerPool: workerPool,
		JobChannel: make(chan BuilderJob),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it. Expects that the formation data for the Plane is preloaded.
func (w Builder) Start(id int) {
	go func() {
		w.ID = id
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				err := w.executeJob(job.Payload)
				if err != nil {
					w.logError(err, "unable to execute job")
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

func (w Builder) logError(err error, msg string) {
	log.Err(err).Stack().Int("WorkerID", w.ID).Msg(msg)
}

func (w Builder) executeJob(payload BuilderPayload) error {
	// we have received a work request.
	err := payload.Plane.Validate()
	if err != nil {
		return fmt.Errorf("invalid plane object: %w", err)
	}
	db, err := db.ConnectToDB()
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	newPlane := Plane{
		Num: payload.Plane.Formation.getNextNum(0),
		FormationID: payload.Plane.FormationID,
	}
	result := db.Save(&newPlane)
	if result.Error != nil {
		return fmt.Errorf("unable to update formation with new planes with error: %w", result.Error)
	}
	return nil
}