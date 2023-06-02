package model

import (
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/db"
	"gorm.io/gorm"
	"log"
)

type Flight struct {
	gorm.Model
	Name string `json:"name"`
}

func SetupFlights() {
	err := db.DBConnection.AutoMigrate(Flight{})
	if err != nil {
		log.Fatalf("Unable to perform database migration: %v", err)
	}
}

func GetAllFlights() []Flight {
	var flights []Flight
	result := db.DBConnection.Find(&flights)
	if result.Error != nil {
		log.Fatalf("Unable to retrieve list of flights: %v", result.Error)
	}
	return flights
}

func GetFlight(id int) Flight {
	var flight Flight
	result := db.DBConnection.First(&flight, id)
	if result.Error != nil {
		log.Printf("Unable to find flight with ID: %d", id)
	}
	return flight
}
