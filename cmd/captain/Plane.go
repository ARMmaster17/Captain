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
	VMID int `validate:"gte=0"`
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
	err = p.buildPlane()
	if err != nil {
		return fmt.Errorf("unable to trigger plane build with error: %w", err)
	}
	return nil
}

func (p *Plane) BeforeDelete(tx *gorm.DB) error {
	// TODO: Glue in old code to delete plane.
	return nil
}

func (p *Plane) Validate() error {
	err := validator.New().Struct(p)
	if err != nil {
		return fmt.Errorf("invalid parameters for plane: %w", err)
	}
	return nil
}

func (p *Plane) buildPlane() error {
	log.Debug().Str("PlaneName", p.getFQDN()).Msg("building new plane")
	px, err := ProxmoxAdapterConnect()
	if err != nil {
		return fmt.Errorf("unable to contact Proxmox cluster with error: %w", err)
	}
	planeConfig, err := ProxmoxBuildLxc(p)
	if err != nil {
		return fmt.Errorf("uanble to build plane configuration with error: %w", err)
	}
	vmid, err := ProxmoxCreateLxc(px, planeConfig)
	if err != nil {
		return fmt.Errorf("unable to build plane with error: %w", err)
	}
	p.VMID = vmid
	return nil
}

func (p *Plane) destroyPlane() error {
	log.Debug().Str("PlaneName", p.getFQDN()).Msg("destroying plane")
	return nil
}

