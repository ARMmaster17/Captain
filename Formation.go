package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Formation struct {
	gorm.Model
	Name string
	CPU int
	RAM int
	BaseName	string
	Domain		string
	TargetCount	int
	Planes []Plane
	FlightID int
	Flight Flight
}

func initFormations(db *gorm.DB) error {
	err := initPlanes(db)
	if err != nil {
		return fmt.Errorf("unable to migrate formation dependencies with error: %w", err)
	}
	err = db.AutoMigrate(&Formation{})
	if err != nil {
		return fmt.Errorf("unable to migrate formation schema with error: %w", err)
	}
	return nil
}
