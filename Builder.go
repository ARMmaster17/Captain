package main

import (
	"fmt"
	"github.com/ARMmaster17/Captain/db"
	"github.com/rs/zerolog/log"
	"sync"
)

type Builder struct {
	ID			int
}

func (w Builder) logError(err error, msg string) {
	log.Err(err).Stack().Int("WorkerID", w.ID).Msg(msg)
}

func (w Builder) buildPlane(payload Plane, wg *sync.WaitGroup) {
	defer wg.Done()
	// we have received a work request.
	err := payload.Validate()
	if err != nil {
		w.logError(err, fmt.Sprintf("Invalid plane object"))
		return
	}
	db, err := db.ConnectToDB()
	if err != nil {
		w.logError(err, fmt.Sprintf("unable to connect to database"))
		return
	}
	newPlane := Plane{
		Num: payload.Num,
		FormationID: payload.FormationID,
	}
	result := db.Save(&newPlane)
	if result.Error != nil {
		w.logError(err, fmt.Sprintf("unable to update formation with new planes"))
		return
	}
}
