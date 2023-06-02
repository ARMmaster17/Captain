package v1

import (
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/api/v1/auth"
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/api/v1/srsx"
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/api/v1/util"
	"github.com/gin-gonic/gin"
)

func GenerateRoutes(group *gin.RouterGroup) {
	v1Group := group.Group("/v1")
	v1Group.POST("/util/ping", util.HandlePingPost)
	v1Group.POST("/auth/login", auth.HandleLoginPost)
	v1Group.POST("/auth/refresh", auth.HandleRefreshPost)
	v1Group.GET("/srsx/flights", srsx.HandleFlights)
	v1Group.POST("/srsx/flight", srsx.HandleCreateFlight)
	v1Group.GET("/srsx/flight/:id", srsx.HandleReadFlight)
}
