package main

import (
	v1 "github.com/ARMmaster17/Captain/airspace_controller/pkg/api/v1"
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/middleware"
	"github.com/ARMmaster17/Captain/airspace_controller/pkg/model"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Setup
	model.SetupAll()

	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.Use(static.Serve("/", static.LocalFile("./frontend/dist", true)))
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.File("./frontend/dist/index.html")
		} else {
			c.String(http.StatusNotFound, "")
		}
	})

	// Build routes
	apiRoutes := router.Group("/api")
	v1.GenerateRoutes(apiRoutes)

	// Run
	log.Fatal(router.Run(":3000"))
}
