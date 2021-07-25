// +build windows

package framework

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB connects to a supported database using the provided connection string.
func InitDB(connectionString string) (*gorm.DB, error) {

	switch connectionString {
	case "":
		return gorm.Open(postgres.Open("host=localhost user=captain password=captain dbname=captain port=9920 sslmode=disable"), &gorm.Config{})
	case "TEST":
		return nil, nil
	default:
		return gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	}
}
