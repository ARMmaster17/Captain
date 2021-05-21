package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerAirspaceHandlers(router *gin.Engine) {
	router.GET(airspacePath(""), handleAirspaceAllGet)
	router.POST(airspacePath(""), handleAirspaceNewPost)
	router.POST(airspacePath("/:airspaceid/delete"), handleAirspaceDelete)
}

func airspacePath(uri string) string {
	return fmt.Sprintf("/airspace%s", uri)
}

func handleAirspaceAllGet(c *gin.Context) {
	client := getCaptainClient()
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
	client := getCaptainClient()
	// TODO: Validate input.
	_, err := client.CreateAirspace(c.PostForm("HumanName"), c.PostForm("NetName"))
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	} else {
		c.Redirect(http.StatusFound, "/airspace")
	}
}

func handleAirspaceDelete(c *gin.Context) {
	client := getCaptainClient()
	airspaceID, err := getUrlIDParameter("airsapce", c)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	err = client.DeleteAirspace(airspaceID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	c.Redirect(http.StatusFound, "/airspace")
}
