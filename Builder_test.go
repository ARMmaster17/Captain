package main

import (
	db2 "github.com/ARMmaster17/Captain/db"
	"github.com/go-playground/assert"
	"sync"
	"testing"
)

func TestBuilderCreateDestroyCycle(t *testing.T) {
	db, err := db2.ConnectToDB()
	assert.Equal(t, err, nil)

	builder := Builder{
		ID: 1,
	}
	flight := Flight{
		Name:       "Test Flight",
		AirspaceID: 0,
	}
	db.Create(&flight)
	formation := Formation{
		Name:        "TestFormation",
		CPU:         1,
		RAM:         128,
		Disk:        4,
		BaseName:    "testformation",
		Domain:      "example.com",
		TargetCount: 0,
		FlightID:    int(flight.ID),
	}
	db.Create(&formation)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	builder.buildPlane(Plane{
		Formation: formation,
		FormationID: int(formation.ID),
		Num: formation.getNextNum(0),
	}, wg)
	// TODO: Check that plane got built right.
	db.Where("formation_id = ?", &formation.ID).Delete(Plane{})
	db.Delete(Formation{}, &formation.ID)
	db.Delete(Flight{}, &flight.ID)
}