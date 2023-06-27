package model

import (
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/db"
	"gorm.io/gorm"
	"log"
)

type Formation struct {
	gorm.Model
	FlightID    int    `json:"flight_id"`
	Name        string `json:"name"`
	CpuCount    int    `json:"cpu_count"`
	RamMB       int    `json:"ram_mb"`
	DiskGB      int    `json:"disk_gb"`
	TargetScale int    `json:"target_scale"`
}

func SetupFormations() {
	err := db.DBConnection.AutoMigrate(Formation{})
	if err != nil {
		log.Fatalf("Unable to perform database migration: %v", err)
	}
}
