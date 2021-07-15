package models

import "github.com/go-playground/validator"

// CRUDObject is a generic object that is stored semi-permanently in a configured SQL-based database.
type CRUDObject interface {
	Create() error
	GetByID() (CRUDObject, error)
	GetAll() ([]CRUDObject, error)
	Update() error
	Delete() error
}

// Validate checks all fields of the given object using the validator library. Errors will be returned for
// each invalid value in the struct.
func Validate(c CRUDObject) error {
	return validator.New().Struct(c)
}
