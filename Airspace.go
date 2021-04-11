package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Airspace struct {
	gorm.Model
	HumanName string
	NetName string
	Flights []Flight
}

func initAirspaces(db *gorm.DB) error {
	log.Debug().Msg("initializing flights")
	err := initFlights(db)
	if err != nil {
		return fmt.Errorf("unable to migrate airspace schema dependencies with error: %w", err)
	}
	log.Debug().Msg("performing airspace schema migrations")
	err = db.AutoMigrate(&Airspace{})
	if err != nil {
		return fmt.Errorf("unable to migrate airspace schema with error: %w", err)
	}
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
			return fmt.Errorf("unable to create default airspace with error: %w", result.Error)
		}
	}
	return nil
}

func (a *Airspace) performHealthChecks(db *gorm.DB) error {
	result := db.Where("airspace_id = ?", a.ID).Find(&a.Flights)
	if result.Error != nil {
		return fmt.Errorf("unable to list flights for airspace %s with error: %w", a.HumanName, result.Error)
	}
	for i := 0; i < len(a.Flights); i++ {
		log.Trace().Str("airspace", a.NetName).Str("flight", a.Flights[i].Name).Msg("checking health of flight")
		err := a.Flights[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to check health of flight %s with error: %w", a.Flights[i].Name, err)
		}
	}
	return nil
}
