package main

import "testing"

func TestConnectToDB(t *testing.T) {
	db, err := ConnectToDB()
	if err != nil {
		t.Errorf("unexpected error in creating Sqlite3 database: %w", err)
	}
	if db == nil {
		t.Errorf("db unexpectedly nil")
	}
}
