package db

import (
	"os"
	"testing"
)

func TestConnectToDB(t *testing.T) {
	db, err := ConnectToDB()
	if err != nil {
		t.Errorf("unexpected error in creating Sqlite3 database: %w", err)
	}
	if db == nil {
		t.Errorf("db unexpectedly nil")
	}
}

func TestDBDefaultConnectionString(t *testing.T) {
	connString, err := getDBConnectionString()
	if err != nil {
		t.Errorf("unexpected error getting default connection string: %w", err)
		return
	}
	if connString != os.Getenv("CAPTAIN_DB") {
		t.Errorf("expected connection string '%s', got %s", os.Getenv("CAPTAIN_DB"), connString)
	}
}

func TestGetConfiguredDBDriverEnv(t *testing.T) {
	db, err := getConfiguredDBDriver()
	if err != nil {
		t.Errorf("unexpected error connnecting to database: %w", err)
	}
	if db.Name() != "sqlite" {
		t.Errorf("expected DB adapter of type 'sqlite', got '%s'", db.Name())
	}
}
