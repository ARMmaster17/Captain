package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DBConnection = dbSetup()

func dbSetup() *gorm.DB {
	// Setup
	db, err := gorm.Open(postgres.Open(os.Getenv("CAPTAIN_DATABASE")))
	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	// Done
	return db
}
