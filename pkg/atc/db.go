package atc

import (
	"database/sql"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"fmt"
)

func getDBConnection() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./captain.db")
	if err != nil {
		return nil, fmt.Errorf("unable to initialize database with error: %w", err)
	}
	return database, nil
}

func prepDB() (*sql.DB, error) {
	database, err := getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to prepare to apply migrations to the database with error: %w", err)
	}
	driver, err := sqlite3.WithInstance(database, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize database driver with error: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"sqlite3", driver)
	m.Steps(2)
	return database, nil
}

func Execute(query string) error {
	database, err := getDBConnection()
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
}
