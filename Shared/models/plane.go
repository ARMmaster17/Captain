package models

import (
	"gorm.io/gorm"
)

type Plane struct{
	gorm.Model
	Name string `validate:"required,alphanum"`
	CPU int `validate:"gt=0,lt=8192"`
	RAM int `validate:"gt=15"`
	Disk int `validate:"gt=0"`
}

func NewPlane() *Plane {
	return &Plane{
		Name:  "default",
		CPU: 1,
		RAM: 512,
		Disk: 8,
	}
}

func (p *Plane) Create() error {
	panic("implement me")
}

func (p *Plane) GetByID() (CRUDObject, error) {
	panic("implement me")
}

func (p *Plane) GetAll() ([]CRUDObject, error) {
	panic("implement me")
}

func (p *Plane) Update() error {
	panic("implement me")
}

func (p *Plane) Delete() error {
	panic("implement me")
}

