package dto

import (
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/model"
)

type FlightDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ConvertFlightModelToDTO(flight model.Flight) FlightDTO {
	return FlightDTO{
		ID:   flight.ID,
		Name: flight.Name,
	}
}

func ConvertFlightDTOToModel(dto FlightDTO, allowId bool) model.Flight {
	flight := model.Flight{
		Name: dto.Name,
	}
	if allowId {
		flight.ID = dto.ID
	}
	return flight
}
