package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Airspace struct {
	gorm.Model
	HumanName string
	NetName string
	Flights []Flight
}

func initAirspaces(db *gorm.DB) error {
	err := initFlights(db)
	if err != nil {
		return fmt.Errorf("unable to migrate airspace schema dependencies with error: %w", err)
	}
	err = db.AutoMigrate(&Airspace{})
	if err != nil {
		return fmt.Errorf("unable to migrate airspace schema with error: %w", err)
	}
	var airspaceCount int64
	db.Model(&Airspace{}).Count(&airspaceCount)
	if airspaceCount == 0 {
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
	for i := 0; i < len(a.Flights); i++ {
		err := a.Flights[i].performHealthChecks(db)
		if err != nil {
			return fmt.Errorf("unable to check health of flight %s with error: %w", a.Flights[i].Name, err)
		}
	}
	return nil
}
