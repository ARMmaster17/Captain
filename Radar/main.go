package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	registerRootHandlers(r)
	registerAirspaceHandlers(r)
	err := r.Run("0.0.0.0:5001")
	if err != nil {
		// TODO: log error
		return
	}
}
