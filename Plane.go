package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Plane struct {
	gorm.Model
	Name string
	VMID string
	FormationID int
	Formation Formation
}

func initPlanes(db *gorm.DB) error {
	err := db.AutoMigrate(&Plane{})
	if err != nil {
		return fmt.Errorf("unable to migrate plane schema with error: %w", err)
	}
	return nil
}
