package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// registerFlightHandlers creates all the flight routes in the given router instance. Requires registerAirspaceHandlers
// to run first to ensure successful bindings.
func registerFlightHandlers(router *gin.Engine) {
	router.GET(flightPath(""), handleFlightAllGet)
	router.POST(flightPath(""), handleFlightNewPost)
	router.POST(flightPath("/:flightid/delete"), handleFlightDelete)
}

// flightPath helper method to create a fully-qualified path to the other flight pages.
func flightPath(uri string) string {
	return fmt.Sprintf("%s/:airspaceid%s", airspacePath(""), uri)
}

// handleFlightAllGet handles requests to the main flight page. A listing of all flights in the parent airspace
// is rendered.
func handleFlightAllGet(c *gin.Context) {
	client := getCaptainClient()
	airspace, err := getAirspaceFromURLParameter(c, client)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	flights, err := client.GetFlightsByAirspace(airspace.ID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	c.HTML(http.StatusOK, "flight/index.html", gin.H{
		"flights":  flights,
		"airspace": airspace,
		"pagename": airspace.HumanName,
	})
}

// handleFlightNewPost handles requests to create a new flight in the parent airspace with the given form parameters.
// If the request is successful, the user will be redirected to the flight listings page in the parent airspace.
func handleFlightNewPost(c *gin.Context) {
	client := getCaptainClient()
	// TODO: Validate input.
	airspace, err := getAirspaceFromURLParameter(c, client)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid airspace ID: %w", err))
		return
	}
	_, err = client.CreateFlight(c.PostForm("Name"), airspace.ID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d", airspace.ID))
	}
}

// handleFlightDelete handles requests to delete a flight in the parent airspace with the given ID. If the request
// is successful, the user will be redirected to the flight listings page in the parent airspace.
func handleFlightDelete(c *gin.Context) {
	client := getCaptainClient()
	airspace, err := getAirspaceFromURLParameter(c, client)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	flight, err := getFlightFromURLParameter(c, client)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	err = client.DeleteFlight(flight.ID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d", airspace.ID))
}
