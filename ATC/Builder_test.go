package main

import (
	"fmt"
	db2 "github.com/ARMmaster17/Captain/ATC/DB"
	"github.com/go-playground/assert"
	"sync"
	"testing"
)

func TestBuilderCreateDestroyCycle(t *testing.T) {
	db, err := db2.ConnectToDB()
	assert.Equal(t, err, nil)

	builder := builder{
		ID: 1,
	}
	flight := Flight{
		Name:       "Test Flight",
		AirspaceID: 1,
	}
	tx := db.Create(&flight)
	assert.Equal(t, tx.Error, nil)
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
	tx = db.Create(&formation)
	assert.Equal(t, tx.Error, err)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	builder.buildPlane(Plane{
		Formation: formation,
		FormationID: int(formation.ID),
		Num: formation.getNextNum(0),
	}, wg)
	fmt.Println(formation.ID)
	// TODO: Check that plane got built right.
	tx = db.Where("formation_id = ?", formation.ID).Delete(&Plane{})
	assert.Equal(t, tx.Error, err)
	tx = db.Delete(&Formation{}, formation.ID)
	assert.Equal(t, tx.Error, err)
	tx = db.Delete(&Flight{}, flight.ID)
	assert.Equal(t, tx.Error, err)
}