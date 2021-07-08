package models

import (
	"gorm.io/gorm"
)

type Plane struct{
	gorm.Model
	Name string `validate:"required,alphanum"`
}

func NewPlane() *Plane {
	return &Plane{
		Name:  "default",
	}
}

func (Plane) Create() error {
	panic("implement me")
}

func (Plane) GetByID() (CRUDObject, error) {
	panic("implement me")
}

func (Plane) GetAll() ([]CRUDObject, error) {
	panic("implement me")
}

func (Plane) Update() error {
	panic("implement me")
}

func (Plane) Delete() error {
	panic("implement me")
}

