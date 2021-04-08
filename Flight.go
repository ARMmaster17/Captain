package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	Name string
	Formations []Formation
	AirspaceID int
	Airspace Airspace
}

func initFlights(db *gorm.DB) error {
	err := initFormations(db)
	if err != nil {
		return fmt.Errorf("unable to migrate flight schema dependencies with error: %w", err)
	}
	err = db.AutoMigrate(&Flight{})
	if err != nil {
		return fmt.Errorf("unable to migrate flight schema with error: %w", err)
	}
	return nil
}
