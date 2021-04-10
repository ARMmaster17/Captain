package main

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {
	// TODO: Have this as a configurable option
	var dbPath = "test.db"
	log.Debug().Str("dbPath", dbPath).Msg("connecting to Sqlite3 database")
	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
