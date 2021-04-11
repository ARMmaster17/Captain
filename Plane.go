package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Plane struct {
	gorm.Model
	Num int `validate:"required,gte=1"`
	ProxmoxIdentifier int `validate:"gte=0"`
	FormationID int
	Formation Formation `validate:"-"`
}

func NewPlane(num int) (*Plane, error) {
	plane := Plane{
		Num: num,
	}
	err := plane.Validate()
	if err != nil {
		return nil, fmt.Errorf("unable to create plane with error: %w", err)
	}
	return &plane, nil
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
	err := p.Validate()
	if err != nil {
		return fmt.Errorf("unable to create plane: %w", err)
	}
	result := tx.First(p.Formation, p.FormationID)
	if result.Error != nil {
		return fmt.Errorf("unable to get formation ID %d for new plane: %w", p.FormationID, err)
	}
	err = p.Formation.Validate()
	if err != nil {
		return fmt.Errorf("invalid formation configuration for plane: %w", err)
	}
	err = p.buildPlane(tx)
	if err != nil {
		// TODO: Should detect what kind of error occurred
		return fmt.Errorf("unable to trigger plane build with error: %w", err)
	}
	return nil
}

func (p *Plane) BeforeDelete(tx *gorm.DB) error {
	return p.destroyPlane()
}

func (p *Plane) Validate() error {
	err := validator.New().Struct(p)
	if err != nil {
		return fmt.Errorf("invalid parameters for plane: %w", err)
	}
	return nil
}

func (p *Plane) buildPlane(db *gorm.DB) error {
	log.Debug().Str("PlaneName", p.getFQDN()).Msg("building new plane")
	px, err := ProxmoxAdapterConnect()
	if err != nil {
		return fmt.Errorf("unable to contact Proxmox cluster with error: %w", err)
	}
	err = ProxmoxBuildLxc(db, px, p)
	if err != nil {
		return fmt.Errorf("unable to build plane: %w", err)
	}
	return nil
}

func (p *Plane) destroyPlane() error {
	log.Debug().Str("PlaneName", p.getFQDN()).Msg("destroying plane")
	px, err := ProxmoxAdapterConnect()
	if err != nil {
		return fmt.Errorf("unable to connect to Proxmox cluster: %w", err)
	}
	err = ProxmoxDestroyLxc(px, p)
	if err != nil {
		return fmt.Errorf("unable to destroy plane %s: %w", p.getFQDN(), err)
	}
	return nil
}

