package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerFlightHandlers(router *gin.Engine) {
	router.GET(flightPath(""), handleFlightAllGet)
	router.POST(flightPath(""), handleFlightNewPost)
	router.POST(flightPath("/:flightid/delete"), handleFlightDelete)
}

func flightPath(uri string) string {
	return fmt.Sprintf("%s/:airspaceid%s", airspacePath(""), uri)
}

func handleFlightAllGet(c *gin.Context) {
	client := getCaptainClient()
	airspaceID, err := getUrlIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid airspace ID: %w", err))
		return
	}
	airspace, err := client.GetAirspaceByID(airspaceID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	flights, err := client.GetFlightsByAirspace(airspaceID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	c.HTML(http.StatusOK, "flight/index.html", gin.H{
		"flights":  flights,
		"airspace": airspace,
	})
}

func handleFlightNewPost(c *gin.Context) {
	client := getCaptainClient()
	// TODO: Validate input.
	airspaceID, err := getUrlIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid airspace ID: %w", err))
		return
	}
	_, err = client.CreateFlight(c.PostForm("Name"), airspaceID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d", airspaceID))
	}
}

func handleFlightDelete(c *gin.Context) {
	client := getCaptainClient()
	airspaceID, err := getUrlIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	flightID, err := getUrlIDParameter("flight", c)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	err = client.DeleteFlight(flightID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d", airspaceID))
}
