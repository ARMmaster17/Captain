package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := initRouter()
	err := r.Run("0.0.0.0:5001")
	if err != nil {
		// TODO: log error
		return
	}
}

func initRouter() *gin.Engine {
	r := gin.Default()
	executablePath, _ := os.Executable()
	if executablePath == "" {
		r.LoadHTMLGlob("/etc/captain/radar/templates/**/*")
	} else {
		r.LoadHTMLGlob("templates/**/*")
	}
	registerRootHandlers(r)
	registerAirspaceHandlers(r)
	registerFlightHandlers(r)
	registerFormationHandlers(r)
	return r
}
