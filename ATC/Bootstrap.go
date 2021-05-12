package main

import (
	"fmt"
	"github.com/ARMmaster17/Captain/ATC/DB"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Creates the entire Captain stack on target Proxmox cluster. Assumes that the current ATC instance is running in
// standalone mode with a Sqlite3 database. After this procedure runs, the standalone ATC instance will be converted
// into cluster mode using a Postgres database with distributed builders.
func BootstrapCluster() error {
	log.Info().Msg("bootstrapping cluster...")
	log.Debug().Msg("connecting to database")
	db, err := DB.ConnectToDB()
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	log.Debug().Msg("running migrations")
	err = initAirspaces(db)
	if err != nil {
		return fmt.Errorf("unable to perform schema migrations on database: %w", err)
	}
	airspace, err := bootstrapCreateSystemAirspace(db)
	if err != nil {
		return fmt.Errorf("unable to create airspace for Captain services: %w", err)
	}
	err = bootstrapCreateCaptainServices(db, airspace)
	if err != nil {
		return fmt.Errorf("unable to provision Captain stack: %w", err)
	}
	return nil
}

// Handles the actual creation of the Captain Stack. Assumes that the local ATC instance has been fully initialized
// and is ready to start creating formations.
func bootstrapCreateCaptainServices(db *gorm.DB, airspace *Airspace) error {
	airspace.Flights = append(airspace.Flights, Flight{
		Name: "Captain Core Services",
	})
	results := db.Save(&airspace)
	if results.Error != nil {
		return fmt.Errorf("unable to create core services flight: %w", results.Error)
	}
	flight := airspace.Flights[0]
	flight.Formations = append(flight.Formations, Formation{
		Name:        "Captain Server",
		CPU:         1,
		RAM:         256,
		Disk:        8,
		BaseName:    "captain",
		Domain:      "core.cap",
		TargetCount: 1,
	})
	flight.Formations = append(flight.Formations, Formation{
		Name:        "PostgreSQL Server",
		CPU:         1,
		RAM:         256,
		Disk:        12,
		BaseName:    "psql",
		Domain:      "core.cap",
		TargetCount: 1,
	})
	results = db.Save(&flight)
	if results.Error != nil {
		return fmt.Errorf("unable to create core service flights: %w", results.Error)
	}
	err := flight.performHealthChecks(db)
	if err != nil {
		return fmt.Errorf("unable to create core service planes: %w", err)
	}

	return nil
}

// Ensures that this node is clean and does not have any other initialized airspaces. If this node is clean, a managed
// airspace will be created that holds all the internal services used by Captain to avoid collisions with user-created
// flights and formations.
func bootstrapCreateSystemAirspace(db *gorm.DB) (*Airspace, error) {
	airspacesExist, err := bootstrapCheckIfAirspacesAlreadyExist(db)
	if err != nil {
		return nil, fmt.Errorf("unable to check existing airspaces: %w", err)
	}
	if airspacesExist {
		return nil, fmt.Errorf("cluster is not empty, airspaces have already been initialized")
	}
	systemAirspace := Airspace{
		HumanName: "Captain Services",
		NetName:   "captain",
	}
	result := db.Create(&systemAirspace)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to create system airspace")
	}
	return &systemAirspace, nil
}

// Verifies that only the default airspace exists in the state database. Used to verify that the ATC instance is in the
// proper state for bootstrapping.
func bootstrapCheckIfAirspacesAlreadyExist(db *gorm.DB) (bool, error) {
	var airspaces []Airspace
	result := db.Find(&airspaces)
	if result.Error != nil {
		return false, fmt.Errorf("unable to retrieve list of active airspaces: %w", result.Error)
	}
	// TODO: Check names of airspaces to determine that only the blank 'default' namespace exists.
	return result.RowsAffected > 1, nil
}
