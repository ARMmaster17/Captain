package main

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/go-playground/validator"
)

type Plane struct {
	gorm.Model
	Num int `validate:"required,gte=0"`
	VMID string
	FormationID int
	Formation Formation
}

func NewPlane(num int) (*Plane, error) {
	plane := Plane{
		Num: num,
	}
	var validate *validator.Validate
	err := validate.Struct(plane)
	if err != nil {
		return nil, fmt.Errorf("invalid parameters for plane configurations: %w", err)
	}
	return &plane, err
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
