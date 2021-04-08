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

func (f *Formation) performHealthChecks(db *gorm.DB) error {
	// Remove dead planes.
	for i := 0; i < len(f.Planes); i++ {
		// TODO: Run plane health checks.
	}
	// Check that the number of active (or planned) planes equals the target.
	if len(f.Planes) < f.TargetCount {
		var offset int = f.TargetCount - len(f.Planes)
		// TODO: Generate unique names for new planes.
		for i := 0; i < offset; i++ {
			f.Planes = append(f.Planes, Plane{
				Name: fmt.Sprintf("%s%d.%s", f.BaseName, i, f.Domain),
			})
		}
		result := db.Save(f)
		if result.Error != nil {
			return fmt.Errorf("unable to update formation with new planes with error: %w", result.Error)
		}
	}
	return nil
}
