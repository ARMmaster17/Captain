package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Airspace model. Represents an isolated space for applications to run independent of other airspaces.
type Airspace struct {
	gorm.Model
	// A human-friendly name to be used in GUIs and CLI outputs for identification purposes.
	HumanName string
	// A computer friendly name. In the future, this name will be used in the FQDN of any flights in the airspace
	// (e.x. formation1.flight.airspace.example.com).
	NetName string
	Flights []Flight
}

// Performs migrations on the Airspace schema. Will also create a default airspace if no airspaces exist
// in the state database. Will also trigger all dependent migrations as well.
func initAirspaces(db *gorm.DB) error {
	// Initialize dependent schemas first.
	log.Debug().Msg("initializing flights")
	err := initFlights(db)
	if err != nil {
		return fmt.Errorf("unable to migrate airspace schema dependencies with error:\n%w", err)
	}
	log.Debug().Msg("performing airspace schema migrations")
	err = db.AutoMigrate(&Airspace{})
	if err != nil {
		return fmt.Errorf("unable to migrate airspace schema with error:\n%w", err)
	}
	// Create a default airspace if none exist.
	log.Trace().Msg("checking if default airspace exists")
	var airspaceCount int64
	db.Model(&Airspace{}).Count(&airspaceCount)
	if airspaceCount == 0 {
		log.Trace().Msg("default airspace does not exist, creating...")
		airspace := Airspace{
			HumanName: "Default Airspace",
			NetName: "default",
		}
		result := db.Create(&airspace)
		if result.Error != nil {
			return fmt.Errorf("unable to create default airspace with error:\n%w", result.Error)
		}
	}
	return nil
}

// Checks the health of all dependent objects in the airspace.
func (a *Airspace) performHealthChecks(db *gorm.DB) error {
	result := db.Where("airspace_id = ?", a.ID).Find(&a.Flights)
	if result.Error != nil {
		return fmt.Errorf("unable to list flights for airspace %s with error:\n%w", a.HumanName, result.Error)
	}
	for i := 0; i < len(a.Flights); i++ {
		log.Trace().Str("airspace", a.NetName).Str("flight", a.Flights[i].Name).Msg("checking health of flight")
		err := a.Flights[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to check health of flight %s with error:\n%w", a.Flights[i].Name, err)
		}
	}
	return nil
}
