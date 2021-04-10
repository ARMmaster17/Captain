package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	Name string `validate:"required,min=1"`
	Formations []Formation `validate:"-"`
	AirspaceID int
	Airspace Airspace `validate:"-"`
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
	result := db.Where("flight_id = ?", f.ID).Find(&f.Formations)
	if result.Error != nil {
		return fmt.Errorf("unable to list formations for flight %s with error: %w", f.Name, result.Error)
	}
	for i := 0; i < len(f.Formations); i++ {
		log.Trace().Str("flight", f.Name).Str("formation", f.Formations[i].Name).Msg("checking health of formation")
		err := f.Formations[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to perform health check on formation %s with error: %w", f.Formations[i].Name, err)
		}
	}
	return nil
}

func (f *Flight) BeforeCreate(tx *gorm.DB) error {
	err := f.Validate()
	if err != nil {
		return fmt.Errorf("unable to validate flight: %w", err)
	}
	return nil
}

func (f *Flight) Validate() error {
	err := validator.New().Struct(f)
	if err != nil {
		return fmt.Errorf("invalid flight object: %w", err)
	}
	return nil
}
