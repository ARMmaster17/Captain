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

func DBExecuteWithParams(query string, values ...string) error {
	database, err := getDBConnection()
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
	statement, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
	_, err = statement.Exec(values)
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
	err = statement.Close()
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
	err = database.Close()
	if err != nil {
		return fmt.Errorf("unable to execute SQL '%s' with error: %w", query, err)
	}
	return nil
}

func DBQueryWithParams(query string, values ...string) (*sql.Rows, error) {
	database, err := getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to execute SQL query '%s' with error: %w", query, err)
	}
	rows, err := database.Query(query, values)
	if err != nil {
		return nil, fmt.Errorf("unable to execute SQL query '%s' with error: %w", query, err)
	}
	err = database.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to execute SQL query '%s' with error: %w", query, err)
	}
	return rows, nil
}
