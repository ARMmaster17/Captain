package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// A flight is a logical representation of an app. A flight contains several scalable formations that are connected
// together to provide a service. For example, a flight would contain the following formations: web servers, reverse
// proxies, and a database cluster.
type Flight struct {
	gorm.Model
	// Unique name for flight. In the future will be used as part of the FQDN of each plane.
	Name       string      `validate:"required,min=1"`
	Formations []Formation `validate:"-"`
	AirspaceID int
	Airspace   Airspace `validate:"-"`
}

// Performs the necessary migrations to initialize the database to hold flight state information. Triggers all
// needed sub-migrations for all dependencies.
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

// Checks the health of a flight, and all formations within that flight. No real logic is performed here, as the health
// of a flight depends completely on the health of all the dependent formations and their underlying planes.
func (f *Flight) performHealthChecks(db *gorm.DB) error {
	result := db.Where("flight_id = ?", f.ID).Find(&f.Formations)
	if result.Error != nil {
		return fmt.Errorf("unable to list formations for flight %s with error: %w", f.Name, result.Error)
	}
	for i := 0; i < len(f.Formations); i++ {
		log.Trace().Str("flight", f.Name).Str("formation", f.Formations[i].Name).Msg("checking health of formation")
		// TODO: convert to goroutine and waitgroup?
		err := f.Formations[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to perform health check on formation %s with error: %w", f.Formations[i].Name, err)
		}
	}
	return nil
}

// Validates the properties of a model before it is committed into the database.
func (f *Flight) BeforeCreate(tx *gorm.DB) error {
	err := f.Validate()
	if err != nil {
		return fmt.Errorf("unable to validate flight: %w", err)
	}
	return nil
}

// Creates and calls a validator object to verify the validity of the fields of the specified flight object.
func (f *Flight) Validate() error {
	err := validator.New().Struct(f)
	if err != nil {
		return fmt.Errorf("invalid flight object: %w", err)
	}
	return nil
}
