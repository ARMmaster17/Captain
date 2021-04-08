package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Plane struct {
	gorm.Model
	Num int
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

func (p *Plane) isHealthy(db *gorm.DB) (bool, error) {
	// TODO: Actually implement health checks
	return true, nil
}

func (p *Plane) getFQDN() string {
	return fmt.Sprintf("%s%d.%s", p.Formation.BaseName, p.Num, p.Formation.Domain)
}

func (p *Plane) BeforeCreate(tx *gorm.DB) error {
	// TODO: Glue in old code to create plane.
	return nil
}

func (p *Plane) BeforeDelete(tx *gorm.DB) error {
	// TODO: Glue in old code to delete plane.
	return nil
}
