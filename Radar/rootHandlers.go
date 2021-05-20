package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerRootHandlers(router *gin.Engine) {
	router.GET("/", handleRootGet)
}

func handleRootGet(c *gin.Context) {
	c.HTML(http.StatusOK, "root/index.html", gin.H{
		"pagename": "Home",
		"title": "Homepage",
	})
}
