package framework

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB connects to a supported database using the provided connection string.
func InitDB(connectionString string) (*gorm.DB, error) {
	switch connectionString {
	case "":
		log.Warn().Msg("no DB is configured, using in-memory database")
		return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	case "TEST":
		return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	default:
		return gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	}
}
