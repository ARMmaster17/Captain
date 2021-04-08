package main

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"time"
)

func StartMonitoring() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("unable to open database with error: %w", err)
	}
	err = initAirspaces(db)
	if err != nil {
		return fmt.Errorf("unable to migrate database with error: %w", err)
	}
	for {
		err = monitoringLoop(db)
		if err != nil {
			return fmt.Errorf("unable to perform timed checks with error: %w", err)
		}
		time.Sleep(15 * time.Second)
	}
}

func monitoringLoop(db *gorm.DB) error {
	var airspaces []Airspace
	result := db.Find(&airspaces)
	if result.Error != nil {
		return fmt.Errorf("unable to retrieve list of airspaces with error: %w", result.Error)
	}
	for i := 0; i < len(airspaces); i++ {
		err := airspaces[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to perform healthchecks on airspace %s with error: %w", airspaces[i].HumanName, err)
		}
	}
	return nil
}
