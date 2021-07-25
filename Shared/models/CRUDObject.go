package models

import (
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

// CRUDObject is a generic object that is stored semi-permanently in a configured SQL-based database.
type CRUDObject interface {
	Create(db *gorm.DB) error
	GetByID(db *gorm.DB, id int) error
	Update(db *gorm.DB) error
	Delete(db *gorm.DB) error
}

// Validate checks all fields of the given object using the validator library. Errors will be returned for
// each invalid value in the struct.
func Validate(c CRUDObject) error {
	return validator.New().Struct(c)
}
