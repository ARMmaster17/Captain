package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerFormationHandlers(router *gin.Engine) {
	router.GET(formationPath(""), handleFormationAllGet)
	router.POST(formationPath(""), handleFormationNewPost)
	router.POST(formationPath("/:formationid/delete"), handleFormationDelete)
}

func formationPath(uri string) string {
	return fmt.Sprintf("%s/:flightid%s", flightPath(""), uri)
}

func handleFormationAllGet(c *gin.Context) {
	client := getCaptainClient()
	airspaceID, err := getUrlIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid airspace ID: %w", err))
		return
	}
	flightID, err := getUrlIDParameter("flight", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid flight ID: %w", err))
		return
	}
	airspace, err := client.GetAirspaceByID(airspaceID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	flight, err := client.GetFlightByID(flightID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	formations, err := client.GetFormationsByFlight(flightID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	c.HTML(http.StatusOK, "formation/index.html", gin.H{
		"formations": formations,
		"flight":     flight,
		"airspace":   airspace,
	})
}

func handleFormationNewPost(c *gin.Context) {
	client := getCaptainClient()
	// TODO: Validate input.
	airspaceID, err := getUrlIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid airspace ID: %w", err))
		return
	}
	flightID, err := getUrlIDParameter("flight", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid flight ID: %w", err))
		return
	}
	_, err = client.CreateFormation(c.PostForm("Name"),
		flightID,
		forceIntRead(c.PostForm("CPU")),
		forceIntRead(c.PostForm("RAM")),
		forceIntRead(c.PostForm("Disk")),
		c.PostForm("BaseName"),
		c.PostForm("Domain"),
		forceIntRead(c.PostForm("TargetCount")))
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d/%d", airspaceID, flightID))
	}
}

func handleFormationDelete(c *gin.Context) {
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
	formationID, err := getUrlIDParameter("formation", c)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	err = client.DeleteFormation(formationID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error: %w", err))
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d/%d", airspaceID, flightID))
}
