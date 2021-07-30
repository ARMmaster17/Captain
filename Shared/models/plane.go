package models

import (
	"gorm.io/gorm"
)

// Plane is the smallest logical unit in the Captain stack. A plane is synonymous with a VM or container that
// runs a service for a formation/flight.
type Plane struct {
	gorm.Model
	Name string `validate:"required,alphanum"`
	CPU  int    `validate:"gt=0,lt=8192"`
	RAM  int    `validate:"gt=15"`
	Disk int    `validate:"gt=0"`
}

// NewPlane is a factory method for generating Planes with valid defaults.
func NewPlane() *Plane {
	return &Plane{
		Name: "default",
		CPU:  1,
		RAM:  512,
		Disk: 8,
	}
}

// Create commits the given Plane object to the database. Create assumes that the object is new and does not already
// have an assigned ID.
func (p *Plane) Create(db *gorm.DB) error {
	return db.Create(p).Error
}

// GetByID returns a plane with the given integer ID.
func (p *Plane) GetByID(db *gorm.DB, id int) error {
	return db.First(&p, id).Error
}

// Update commits an existing object to the database.
func (p *Plane) Update(db *gorm.DB) error {
	return db.Save(&p).Error
}

// Delete removes a Plane object from the database.
func (p *Plane) Delete(db *gorm.DB) error {
	return db.Delete(&p).Error
}
