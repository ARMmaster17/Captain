package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func StartMonitoring() error {
	db, err := ConnectToDB()
	if err != nil {
		return fmt.Errorf("unable to open database with error: %w", err)
	}
	log.Info().Msg("initializing airspaces")
	err = initAirspaces(db)
	if err != nil {
		return fmt.Errorf("unable to migrate database with error: %w", err)
	}
	log.Info().Msg("beginning monitoring loop on all airspaces")
	for {
		err = monitoringLoop(db)
		if err != nil {
			return fmt.Errorf("unable to perform timed checks with error: %w", err)
		}
		time.Sleep(15 * time.Second)
	}
}

func monitoringLoop(db *gorm.DB) error {
	log.Trace().Msg("retrieving all airspaces from database")
	var airspaces []Airspace
	result := db.Preload(clause.Associations).Find(&airspaces)
	if result.Error != nil {
		return fmt.Errorf("unable to retrieve list of airspaces with error: %w", result.Error)
	}
	for i := 0; i < len(airspaces); i++ {
		log.Trace().Str("airspace", airspaces[i].HumanName).Msg("checking health of airspace")
		err := airspaces[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to perform healthchecks on airspace %s with error: %w", airspaces[i].HumanName, err)
		}
	}
	return nil
}
