package models

import "github.com/go-playground/validator"

type CRUDObject interface {
	Create() error
	GetByID() (CRUDObject, error)
	GetAll() ([]CRUDObject, error)
	Update() error
	Delete() error
}

func Validate(c CRUDObject) error {
	return validator.New().Struct(c)
}
