package atc

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

}

func executeDBCreateCommand(db *sql.DB, query string) error {
	statement, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("unable to prepare create table statement (%s) with error: %w", query, err)
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("unable to create table with error: %w", err)
	}
	err = statement.Close()
	if err != nil {
		return fmt.Errorf("unable to close connection to DB with error: %w", err)
	}
	return nil
}

func RunATC() {
	// Thread 1: checks current environment state against desired state in DB, making corrections as needed

	// Get list of all monitored airspaces
	// for each airspace>flight>formation
	//// Check if healthcheck cooldown has passed
	//// If active checks, perform HTTP ping against all active planes
	//// Check if active flight count does not match target active flight member count
	////// Call Builder() to create/destroy instances as necessary
}

func RunTerminal() {
	// Thread 2: Provides web server and API for access to ATC. (Should probably be moved to another application.
}