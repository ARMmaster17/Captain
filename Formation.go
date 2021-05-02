package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"sync"
)

// Formation is the lowest-level object that is directly addressable by the user (within the context of Captain). A
// formation manages a group of planes that are scaled up and down automatically. A formation is a logical representation
// of an internal service for an application. For example, all web servers that serve the same web app would be
// considered part of the same formation. All planes in a formation will be exactly the same except for the FQDN.
type Formation struct {
	gorm.Model
	// Name of the service. Used only in user-facing queries, and is not used internally. Should be as unique as
	// possible for easy identification.
	Name string `validate:"required,min=1"`
	// Number of CPU cores to assign to each plane. The actual implementation of this varies depending on which
	// provider adapter is used.
	CPU int `validate:"required,gte=1,lte=8192"`
	// Amount of RAM in megabytes to assign to each plane.
	RAM int `validate:"required,gte=1,lte=307200"`
	// Size of disk in gigabytes to assign to each plane. It is important that this disk is big enough to store the
	// container OS in addition to application data.
	Disk int `validate:"required,gte=1"`
	// URL-safe name for each plane in formation. Should be unique within the flight. Will be used in the FQDN of each
	// plane that is provisioned within the formation. For example: formation1.example.com.
	BaseName	string `validate:"alphanum,min=1,max=256"`
	// Domain name that forms the ending of the FQDN for each plane. In the future this will be moved to be the same
	// airspace-wide or stack-wide.
	Domain		string `validate:"required,fqdn,min=1"`
	// Desired number of planes that should be operational at any given moment. At each health check interval,
	// remediations will be made to adjust the number of healthy planes in service until it equals this number.
	TargetCount	int `validate:"gte=0"`
	Planes []Plane `validate:"-"`
	FlightID int
	Flight Flight `validate:"-"`
}

// Performs all needed migrations to store formation state information in the database.
func initFormations(db *gorm.DB) error {
	err := initPlanes(db)
	if err != nil {
		return fmt.Errorf("unable to migrate formation dependencies with error: %w", err)
	}
	err = db.AutoMigrate(&Formation{})
	if err != nil {
		return fmt.Errorf("unable to migrate formation schema with error: %w", err)
	}
	return nil
}

// Checks the health of the specified formation by checking the health of each plane. Remediation will be performed
// until the number of operational planes matches the target number in the state database.
func (f *Formation) performHealthChecks(db *gorm.DB) error {
	result := db.Where("formation_id = ?", f.ID).Preload("Formation").Find(&f.Planes)
	if result.Error != nil {
		return fmt.Errorf("unable to list planes for formation %s with error: %w", f.Name, result.Error)
	}
	// Remove dead planes.
	for i := 0; i < len(f.Planes); i++ {
		log.Trace().Str("formation", f.Name).Str("plane", f.Planes[i].getFQDN()).Msg("checking health of plane")
		isHealthy, err := f.Planes[i].isHealthy(db)
		if err != nil {
			return fmt.Errorf("unable to check health of plane %s with error: %w", f.Planes[i].getFQDN(), err)
		}
		if !isHealthy {
			// TODO: Possibly have a grace period up to X seconds before destroying container?
			result := db.Unscoped().Delete(&f.Planes[i])
			if result.Error != nil {
				return fmt.Errorf("unable to remove unhealthy plane %s with error: %w", f.Planes[i].getFQDN(), result.Error)
			}
		}
	}
	// Check that the number of active (or planned) planes equals the target.
	if len(f.Planes) < f.TargetCount {
		log.Debug().Str("formation", f.Name).Msgf("formation currently has %d planes, expected %d", len(f.Planes), f.TargetCount)
		var offset = f.TargetCount - len(f.Planes)
		wg := new(sync.WaitGroup)
		wg.Add(offset)
		for i := 0; i < offset; i++ {
			f.launchBuilder(i, offset, wg)
		}
		log.Trace().Str("formation", f.Name).Msgf("waiting for %d builder threads to return", offset)
		wg.Wait()
	}

	// Reload plane list in case changes were made
	result = db.Where("formation_id = ?", f.ID).Preload("Formation").Find(&f.Planes)
	if result.Error != nil {
		return fmt.Errorf("unable to list planes for formation %s with error: %w", f.Name, result.Error)
	}

	if len(f.Planes) > f.TargetCount {
		log.Debug().Str("formation", f.Name).Msgf("formation currently has %d planes, expected %d", len(f.Planes), f.TargetCount)
		// Delete oldest planes first (usually the first indexes)
		var numToDelete = len(f.Planes) - f.TargetCount
		for i := 0; i < numToDelete; i++ {
			result := db.Unscoped().Delete(&f.Planes[i])
			if result.Error != nil {
				return fmt.Errorf("unable to delete excess plane %s with error: %w", f.Planes[i].getFQDN(), result.Error)
			}
		}
	}

	return nil
}

func (f *Formation) launchBuilder(id int, totalBuilders int, wg *sync.WaitGroup) {
	builder := Builder{
		ID: id,
	}
	log.Trace().Str("formation", f.Name).Msgf("firing off builder %d/%d to build plane", id, totalBuilders)
	go builder.buildPlane(Plane{
		Formation: *f,
		FormationID: int(f.ID),
		Num: f.getNextNum(id),
	}, wg)
}

// Gets the next available unique ID within a formation.
func (f *Formation) getNextNum(offset int) int {
	var nextNum = 1
	for i := 0; i < len(f.Planes); i++ {
		if f.Planes[i].Num > nextNum {
			nextNum = f.Planes[i].Num + 1
		}
	}
	return nextNum + offset
}

// Validates the attributes of a formation object (Note: does not validate any child or parent objects).
func (f *Formation) Validate() error {
	err := validator.New().Struct(f)
	if err != nil {
		return fmt.Errorf("invalid parameters for formation: %w", err)
	}
	return nil
}

// Handler before creation of a formation object in the state database. Validates the attributes of the formation
// object to be sure that all technical values are url-safe and valid.
func (f *Formation) BeforeCreate(tx *gorm.DB) error {
	err := f.Validate()
	if err != nil {
		return fmt.Errorf("invalid formation object: %w", err)
	}
	return nil
}
