package main

import (
	"fmt"
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
	if connString != ":memory:" {
		t.Errorf("expected connection string ':memory:', got %s", connString)
	}
}

func TestGetDBConnectionStringFile(t *testing.T) {
	dbConn := os.Getenv("DATABASE_CONN")
	if dbConn == "sqlite3-memory" {
		t.Skip()
		return
	}
	err := HelperDBCreateConnectionFile()
	if err != nil {
		t.Errorf("unexpected error creating db.conf: %w", err)
		return
	}
	connString, err := getDBConnectionString()
	if err != nil {
		t.Errorf("unexpected error getting connection string: %w", err)
		return
	}
	if connString != HelperDBGetStringFromTestDatabaseType(os.Getenv("DATABASE_CONN")) {
		t.Errorf("ATC did not read database configuration from db.conf")
	}
	err = HelperDBCleanupConnectionFile()
	if err != nil {
		t.Errorf("unexpected error cleaning up file")
	}
}

func TestGetDBConnectionStringEnv(t *testing.T) {
	dbConn := os.Getenv("DATABASE_CONN")
	if dbConn == "sqlite3-memory" {
		t.Skip()
		return
	}
	HelperDBTestCleanup()
	HelperDBCreateConnectionEnv()
	connString, err := getDBConnectionString()
	if err != nil {
		t.Errorf("unexpected error getting connection string: %w", err)
		return
	}
	if connString != HelperDBGetStringFromTestDatabaseType(os.Getenv("DATABASE_CONN")) {
		t.Errorf("ATC did not read database configuration from db.conf")
	}
	HelperDBTestCleanup()
}

func TestGetConfiguredDBDriverFile(t *testing.T) {
	HelperDBTestCleanup()
	err := HelperDBCreateConnectionFile()
	if err != nil {
		t.Errorf("unexpected error creating connection file: %w", err)
		return
	}
	db, err := getConfiguredDBDriver()
	if err != nil {
		t.Errorf("unexpected error connnecting to database: %w", err)
	}
	switch dbtype := os.Getenv("DATABASE_CONN"); dbtype {
	case "postgres":
		if db.Name() != "postgres" {
			t.Errorf("expected DB adapter of type 'postgres', got '%s'", db.Name())
		}
	case "sqlite3-file":
		if db.Name() != "sqlite" {
			t.Errorf("expected DB adapter of type 'sqlite', got '%s'", db.Name())
		}
	case "sqlite3-memory":
		if db.Name() != "sqlite" {
			t.Errorf("expected DB adapter of type 'sqlite', got '%s'", db.Name())
		}
	default:
		t.Errorf("invalid DATABASE_CONN")
	}
	HelperDBTestCleanup()
}

func TestGetConfiguredDBDriverEnv(t *testing.T) {
	HelperDBTestCleanup()
	HelperDBCreateConnectionEnv()
	db, err := getConfiguredDBDriver()
	if err != nil {
		t.Errorf("unexpected error connnecting to database: %w", err)
	}
	switch dbtype := os.Getenv("DATABASE_CONN"); dbtype {
	case "postgres":
		if db.Name() != "postgres" {
			t.Errorf("expected DB adapter of type 'postgres', got '%s'", db.Name())
		}
	case "sqlite3-file":
		if db.Name() != "sqlite" {
			t.Errorf("expected DB adapter of type 'sqlite', got '%s'", db.Name())
		}
	case "sqlite3-memory":
		if db.Name() != "sqlite" {
			t.Errorf("expected DB adapter of type 'sqlite', got '%s'", db.Name())
		}
	default:
		t.Errorf("invalid DATABASE_CONN")
	}
	HelperDBTestCleanup()
}

func HelperDBTestCleanup() {
	HelperDBCleanupConnectionEnv()
	_ = HelperDBCleanupConnectionFile()
}

// HELPER FUNCTIONS

func HelperDBCreateConnectionEnv() {
	_ = os.Setenv("CAPTAIN_DB", HelperDBGetStringFromTestDatabaseType(os.Getenv("DATABASE_CONN")))
}

func HelperDBCleanupConnectionEnv() {
	_ = os.Setenv("CAPTAIN_DB", "")
}

func HelperDBCleanupConnectionFile() error {
	err := os.Remove("/etc/captain/db.conf")
	if err != nil {
		return fmt.Errorf("unable to delete file db.conf: %w", err)
	}
	return nil
}

func HelperDBCreateConnectionFile() error {

	f, err := os.Create("/etc/captain/db.conf")
	if err != nil {
		return fmt.Errorf("unable to open db.conf: %w", err)
	}

	_, err = f.WriteString(HelperDBGetStringFromTestDatabaseType(os.Getenv("DATABASE_CONN")))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("unable to close writer for db.conf: %w", err)
	}

	return nil
}

func HelperDBGetStringFromTestDatabaseType(dbs string) string {
	switch db := dbs; db {
	case "postgres":
		return "postgres://localhost"
	case "sqlite3-file":
		return "test.db"
	case "sqlite3-memory":
		return ":memory"
	default:
		return ""
	}
}
