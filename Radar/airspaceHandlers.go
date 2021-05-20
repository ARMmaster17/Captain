package main

import (
	"fmt"
	"github.com/ARMmaster17/Captain/CaptainLib"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func registerAirspaceHandlers(router *gin.Engine) {
	router.GET(airspacePath(""), handleAirspaceAllGet)
	router.POST(airspacePath(""), handleAirspaceNewPost)
	router.DELETE(airspacePath("/:airspaceid"), handleAirspaceDelete)
}

func airspacePath(uri string) string {
	return fmt.Sprintf("/airspace%s", uri)
}

func handleAirspaceAllGet(c *gin.Context) {
	client := CaptainLib.NewCaptainClient("http://192.168.1.224:5000/")
	airspaces, err := client.GetAllAirspaces()
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	} else {
		c.HTML(http.StatusOK, "airspace/index.html", gin.H{
			"airspaces": airspaces,
		})
	}
}

func handleAirspaceNewPost(c *gin.Context) {
	// TODO: Create client factory to reduce code duplication.
	client := CaptainLib.NewCaptainClient("http://192.168.1.224:5000/")
	// TODO: Validate input.
	_, err := client.CreateAirspace(c.PostForm("HumanName"), c.PostForm("NetName"))
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	} else {
		c.Redirect(http.StatusFound, "/airspace")
	}
}

func handleAirspaceDelete(c *gin.Context) {
	client := CaptainLib.NewCaptainClient("http://192.168.1.224:5000/")
	airspaceID, err := strconv.ParseInt(c.Param("airspaceid"), 10, 64)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	}
	err = client.DeleteAirspace(int(airspaceID))
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	} else {
		c.Redirect(http.StatusFound, "/airspace")
	}
}
