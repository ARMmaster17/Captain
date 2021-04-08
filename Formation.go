package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"github.com/go-playground/validator"
)

type Formation struct {
	gorm.Model
	Name string `validate:"required,min=1"`
	CPU int `validate:"required,gte=1,lte=8192"`
	RAM int `validate:"required,gte=1,lte=307200"`
	BaseName	string `validate:"alphanum,min=1,max=256"`
	Domain		string `validate:"required,fqdn,min=1"`
	TargetCount	int `validate:"required,gte=0"`
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
	result := db.Where("formation_id = ?", f.ID).Preload("Formation").Find(&f.Planes)
	if result.Error != nil {
		return fmt.Errorf("unable to list planes for formation %s with error: %w", f.Name, result.Error)
	}
	// Remove dead planes.
	for i := 0; i < len(f.Planes); i++ {
		log.Trace().Str("formation", f.Name).Str("plane", f.Planes[i].getFQDN()).Msg("checking health of plane")
		isHealthy, err := f.Planes[i].isHealthy(db)
		if err != nil {
			return fmt.Errorf("unable to check health of plane %s with error: %w", f.Planes[i].getFQDN(), err)
		}
		if !isHealthy {
			// TODO: Possibly have a grace period up to X seconds before destroying container?
			result := db.Unscoped().Delete(&f.Planes[i])
			if result.Error != nil {
				return fmt.Errorf("unable to remove unhealthy plane %s with error: %w", f.Planes[i].getFQDN(), result.Error)
			}
		}
	}
	// Check that the number of active (or planned) planes equals the target.
	if len(f.Planes) < f.TargetCount {
		log.Debug().Str("formation", f.Name).Msgf("formation currently has %d planes, expected %d", len(f.Planes), f.TargetCount)
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

	// Reload plane list in case changes were made
	result = db.Where("formation_id = ?", f.ID).Preload("Formation").Find(&f.Planes)
	if result.Error != nil {
		return fmt.Errorf("unable to list planes for formation %s with error: %w", f.Name, result.Error)
	}

	if len(f.Planes) > f.TargetCount {
		log.Debug().Str("formation", f.Name).Msgf("formation currently has %d planes, expected %d", len(f.Planes), f.TargetCount)
		// Delete oldest planes first (usually the first indexes)
		var numToDelete = len(f.Planes) - f.TargetCount
		for i := 0; i < numToDelete; i++ {
			result := db.Unscoped().Delete(&f.Planes[i])
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

func (f *Formation) Validate() error {
	err := validator.New().Struct(f)
	if err != nil {
		return fmt.Errorf("invalid parameters for formation: %w", err)
	}
	return nil
}
