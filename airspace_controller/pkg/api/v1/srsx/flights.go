package srsx

import (
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/db"
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/dto"
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HandleFlights(c *gin.Context) {
	flights := model.GetAllFlights()
	c.JSON(http.StatusOK, flights)
}

func HandleCreateFlight(c *gin.Context) {
	var flightDto dto.FlightDTO
	err := c.ShouldBindJSON(&flightDto)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	mdf := dto.ConvertFlightDTOToModel(flightDto, false)
	result := db.DBConnection.Create(&mdf)
	if result.Error != nil {
		c.String(http.StatusInternalServerError, result.Error.Error())
	}
	c.String(http.StatusOK, "")
}

func HandleReadFlight(c *gin.Context) {
	flightId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	flight := model.GetFlight(flightId)
	c.JSON(http.StatusOK, dto.ConvertFlightModelToDTO(flight))
}
