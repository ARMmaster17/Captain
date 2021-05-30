package main

import (
	"fmt"
	"github.com/ARMmaster17/Captain/ATC/DB"
	"github.com/ARMmaster17/Captain/ATC/IPAM"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"sync"
)

// Builder instance that handles the creation and destruction of planes. Thread-safe and uses a WaitGroup
// to properly lock resources and not overwhelm the provider API.
type builder struct {
	ID			int
}

func (w builder) logError(err error, msg string) {
	log.Err(err).Stack().Int("WorkerID", w.ID).Msg(msg)
}

func (w builder) buildPlane(payload Plane, wg *sync.WaitGroup, mx *sync.Mutex) {
	defer wg.Done()
	// we have received a work request.
	err := payload.Validate()
	if err != nil {
		w.logError(err, fmt.Sprintf("Invalid plane object"))
		return
	}
	db, err := DB.ConnectToDB()
	if err != nil {
		w.logError(err, fmt.Sprintf("unable to connect to database"))
		return
	}
	newNetID, err := getNewNetworkConfig(db, mx)
	if err != nil {
		w.logError(err, fmt.Sprintf("unable to reserve IP address"))
	}
	newPlane := Plane{
		Num: payload.Num,
		FormationID: payload.FormationID,
		NetID: newNetID,
	}
	result := db.Save(&newPlane)
	if result.Error != nil {
		w.logError(err, fmt.Sprintf("unable to update formation with new planes"))
		return
	}
}

func getNewNetworkConfig(db *gorm.DB, mx *sync.Mutex) (string, error) {
	localIPAM := IPAM.NewIPAM(mx, db)
	newIP, err := localIPAM.GetNewAddress()
	if err != nil {
		return "", err
	}
	return newIP.String(), nil
}
