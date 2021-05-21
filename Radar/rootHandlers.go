package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// registerRootHandlers creates all the routes for semi-static pages in the given router instance.
func registerRootHandlers(router *gin.Engine) {
	router.GET("/", handleRootGet)
}

// handleRootGet handles requests for the web GUI home page.
func handleRootGet(c *gin.Context) {
	c.HTML(http.StatusOK, "root/index.html", gin.H{
		"pagename": "Home",
		"title":    "Homepage",
	})
}
