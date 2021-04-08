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

func (f *Flight) performHealthChecks(db *gorm.DB) error {
	for i := 0; i < len(f.Formations); i++ {
		err := f.Formations[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to perform health check on formation %s with error: %w", f.Formations[i].Name, err)
		}
	}
	return nil
}
