package srsx

import (
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

}

func HandleReadFlight(c *gin.Context) {
	flightId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	flight := model.GetFlight(flightId)
	c.JSON(http.StatusOK, flight)
}
