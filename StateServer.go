package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Starts a monitoring server. Ideally this should be run on it's own thread. This server will, at the specified
// interval, check the desired state in the database, and compare it to the actual state as reported by the underlying
// provider drivers and what is reported by the state database after all health checks are completed.
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

// This function is called once per monitoring interval. Checks that each object represented in the state database is
// healthy, and within the requested operating parameters. If not, the object is modified until the issue is mitigated.
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
