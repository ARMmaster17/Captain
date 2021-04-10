package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"regexp"
)

func ConnectToDB() (*gorm.DB, error) {
	log.Debug().Msg("connecting to database")
	dialector, err := getConfiguredDBDriver()
	if err != nil {
		return nil, fmt.Errorf("unable to detect db driver type: %w", err)
	}
	return gorm.Open(dialector, &gorm.Config{})
}

func getConfiguredDBDriver() (gorm.Dialector, error) {
	dbString, err := getDBConnectionString()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve database configuration: %w", err)
	}
	psqlMatch, _ := regexp.Match("$postgresql://", []byte(dbString))
	if psqlMatch {
		return postgres.Open(dbString), nil
	}
	return sqlite.Open(dbString), nil
}

func getDBConnectionString() (string, error) {
	// Check for db.conf file.
	_, err := os.Stat("/etc/captain/db.conf")
	if os.IsNotExist(err) {
		log.Debug().Msg("configuration file db.conf not found")
	} else {
		content, err := ioutil.ReadFile("/etc/captain/db.conf")
		if err != nil {
			return "", fmt.Errorf("unable to read db.conf: %w", err)
		}
		if string(content) == "" {
			return "", fmt.Errorf("db.conf is empty")
		}
		return string(content), nil
	}

	// Check the environment variable CAPTAIN_DB
	envVarString := os.Getenv("CAPTAIN_DB")
	if envVarString != "" {
		return envVarString, nil
	}

	// Use a default Sqlite3 database path
	return ":memory:", nil
}