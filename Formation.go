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
		isHealthy, err := f.Planes[i].isHealthy(db)
		if err != nil {
			return fmt.Errorf("unable to check health of plane %s with error: %w", f.Planes[i].getFQDN(), err)
		}
		if !isHealthy {
			// TODO: Possibly have a grace period up to X seconds before destroying container?
			result := db.Delete(f.Planes[i])
			if result.Error != nil {
				return fmt.Errorf("unable to remove unhealthy plane %s with error: %w", f.Planes[i].getFQDN(), result.Error)
			}
		}
	}
	// Check that the number of active (or planned) planes equals the target.
	if len(f.Planes) < f.TargetCount {
		var offset = f.TargetCount - len(f.Planes)
		// TODO: Generate unique names for new planes.
		for i := 0; i < offset; i++ {
			f.Planes = append(f.Planes, Plane{
				Num: f.getNextNum(i),
			})
		}
		result := db.Save(f)
		if result.Error != nil {
			return fmt.Errorf("unable to update formation with new planes with error: %w", result.Error)
		}
	}
	if len(f.Planes) > f.TargetCount {
		// Delete oldest planes first (usually the first indexes)
		var numToDelete = len(f.Planes) - f.TargetCount
		for i := 0; i < numToDelete; i++ {
			result := db.Delete(f.Planes[i])
			if result.Error != nil {
				return fmt.Errorf("unable to delete excess plane %s with error: %w", f.Planes[i].getFQDN(), result.Error)
			}
		}
	}

	return nil
}

func (f *Formation) getNextNum(offset int) int {
	var nextNum = 1
	for i := 0; i < len(f.Planes); i++ {
		if f.Planes[i].Num > nextNum {
			nextNum = f.Planes[i].Num + 1
		}
	}
	return nextNum + offset
}
