package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// registerFormationHandlers creates all the formation routes in the given router instance. Requires that both
// registerAirspaceHandlers and registerFlightHandlers have run first to ensure successful bindings.
func registerFormationHandlers(router *gin.Engine) {
	router.GET(formationPath(""), handleFormationAllGet)
	router.POST(formationPath(""), handleFormationNewPost)
	router.POST(formationPath("/:formationid/delete"), handleFormationDelete)
}

// formationPath helper method to create a fully-qualified path to the other formation pages.
func formationPath(uri string) string {
	return fmt.Sprintf("%s/:flightid%s", flightPath(""), uri)
}

// handleFormationAllGet handles requests to the main formation page. A listing of all formations in the parent flight
// is rendered.
func handleFormationAllGet(c *gin.Context) {
	client := getCaptainClient()
	airspace, err := getAirspaceFromURLParameter(c, client)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	flight, err := getFlightFromURLParameter(c, client)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	formations, err := client.GetFormationsByFlight(flight.ID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	c.HTML(http.StatusOK, "formation/index.html", gin.H{
		"formations": formations,
		"flight":     flight,
		"airspace":   airspace,
		"pagename": flight.Name,
	})
}

// handleFormationNewPost handles requests to create a new formation in the parent flight with the given form
// parameters. If the request is successful, the user will be redirected to the formation listings page in the parent
// flight.
func handleFormationNewPost(c *gin.Context) {
	client := getCaptainClient()
	// TODO: Validate input.
	airspaceID, err := getURLIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid airspace ID:\n%w", err))
		return
	}
	flightID, err := getURLIDParameter("flight", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid flight ID:\n%w", err))
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
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d/%d", airspaceID, flightID))
	}
}

// handleFormationDelete handles requests to delete a formation in the parent flight with the given ID. If the request
// is successful, the user will be redirected to the formation listings page in the parent flight.
func handleFormationDelete(c *gin.Context) {
	client := getCaptainClient()
	airspaceID, err := getURLIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	flightID, err := getURLIDParameter("flight", c)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	formationID, err := getURLIDParameter("formation", c)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	err = client.DeleteFormation(formationID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/airspace/%d/%d", airspaceID, flightID))
}
