package main

import (
	"testing"
)

func TestE2EMonitoringLoop(t *testing.T) {
	db, err := ConnectToDB()
	if err != nil {
		t.Errorf("unable to open database with error: %w", err)
		return
	}
	err = initAirspaces(db)
	if err != nil {
		t.Errorf("unable to migrate database with error: %w", err)
		return
	}
	err = monitoringLoop(db)
	if err != nil {
		t.Errorf("unable to perform monitoring with error: %w", err)
	}
}

// TODO: Benchmark monitoringLoop with x1000 instances of localhost (or something like that).
