package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func BootstrapCluster() error {
	log.Info().Msg("bootstrapping cluster...")
	log.Trace().Msg("connecting to database")
	db, err := ConnectToDB()
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	err = initAirspaces(db)
	if err != nil {
		return fmt.Errorf("unable to perform schema migrations on database: %w", err)
	}
	var airspaces []Airspace
	result := db.Find(&airspaces)
	if result.Error != nil {
		return fmt.Errorf("unable to retrieve list of active airspaces: %w", err)
	}
	if result.RowsAffected > 1 {
		return fmt.Errorf("database already has initialized airspaces")
	}
	systemAirspace := Airspace{
		HumanName: "Captain Services",
		NetName:   "captain",
	}
	result = db.Create(&systemAirspace)
	if result.Error != nil {
		return fmt.Errorf("unable to create system airspace")
	}
}
