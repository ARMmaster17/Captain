package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"os"
)

// Represents a running instance on any provider, such as an LXC container, a VM, or a Docker/Kubernetes container
// instance. Planes are usually not modified directly by the API, as they are automatically managed by the built-in
// health checks for each formation. In the event of a configuration change at the formation level, planes are usually
// destroyed and recreated in favor of modifying configuration through a complex provider adapter.
type Plane struct {
	gorm.Model
	Num int `validate:"required,gte=1"`
	ProxmoxIdentifier int `validate:"gte=0"`
	FormationID int
	Formation Formation `validate:"-"`
}

// Creates a temporary plane structure with a supplied formation-specific unique ID. It is assumed that this plane
// struct will be populated later with values pulled from the active provider adapter.
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

// Initializes the plane schema in the state database, performing any migrations necessary.
func initPlanes(db *gorm.DB) error {
	err := db.AutoMigrate(&Plane{})
	if err != nil {
		return fmt.Errorf("unable to migrate plane schema with error: %w", err)
	}
	return nil
}

// Not implemented. Will contact the provider adapter for best method to check that plane exists, and will contact over
// TCP or SSH to verify that the plane is operational. If this method returns false, the plane will be destroyed and
// recreated.
func (p *Plane) isHealthy(db *gorm.DB) (bool, error) {
	// TODO: Actually implement health checks
	return true, nil
}

// Builds the plane's FQDN based on the formation and plane configuration. Used in the local network for SSH, TCP, and
// referencing within the provider driver API.
func (p *Plane) getFQDN() string {
	return fmt.Sprintf("%s%d.%s", p.Formation.BaseName, p.Num, p.Formation.Domain)
}

// Validates the parameters of the plane and triggers a provider driver provisioning pipeline to create the plane as
// an actual container/VM instance. Injects the values from the provider driver in to the plane structure before
// committing the create to the database.
func (p *Plane) BeforeCreate(tx *gorm.DB) error {
	err := p.Validate()
	if err != nil {
		return fmt.Errorf("unable to create plane: %w", err)
	}
	result := tx.First(&p.Formation, p.FormationID)
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

// Triggers a destroy operation through the provider driver to remove the running instance before removing the plane
// from the state database.
func (p *Plane) BeforeDelete(tx *gorm.DB) error {
	return p.destroyPlane()
}

// Validates that all the parameters are within valid constraints and that any values used in the FQDN are url-safe.
func (p *Plane) Validate() error {
	err := validator.New().Struct(p)
	if err != nil {
		return fmt.Errorf("invalid parameters for plane: %w", err)
	}
	return nil
}

// Connects to the proper provider driver to create the plane using the specified parameters. State attribute injection
// is done by the underlying provider driver.
func (p *Plane) buildPlane(db *gorm.DB) error {
	log.Debug().Str("PlaneName", p.getFQDN()).Msg("building new plane")
	fmt.Printf("Value of CAPTAIN_DRY_RUN: '%s'\n", os.Getenv("CAPTAIN_DRY_RUN"))
	if os.Getenv("CAPTAIN_DRY_RUN") != "" {
		return nil
	}
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

// Calls the proper provider driver to delete a running instance. Cleans up any extra networking, SDN, or firewall
// rules that may be attached to the running instance (handled by the underlying provider driver).
func (p *Plane) destroyPlane() error {
	log.Debug().Str("PlaneName", p.getFQDN()).Msg("destroying plane")
	fmt.Printf("Value of CAPTAIN_DRY_RUN: '%s'\n", os.Getenv("CAPTAIN_DRY_RUN"))
	if os.Getenv("CAPTAIN_DRY_RUN") != "" {
		return nil
	}
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

