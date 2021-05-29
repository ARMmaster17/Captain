package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// registerAirspaceHandlers creates all the airspace routes in the given router instance.
func registerAirspaceHandlers(router *gin.Engine) {
	router.GET(airspacePath(""), handleAirspaceAllGet)
	router.POST(airspacePath(""), handleAirspaceNewPost)
	router.POST(airspacePath("/:airspaceid/delete"), handleAirspaceDelete)
}

// airspacePath helper method to create a fully-qualified path to other airspace pages.
func airspacePath(uri string) string {
	return fmt.Sprintf("/airspace%s", uri)
}

// handleAirspaceAllGet handles requests to the main airspace page. A listing of all airspaces is rendered.
func handleAirspaceAllGet(c *gin.Context) {
	client := getCaptainClient()
	airspaces, err := client.GetAllAirspaces()
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
	} else {
		c.HTML(http.StatusOK, "airspace/index.html", gin.H{
			"airspaces": airspaces,
			"pagename": "Airspaces",
		})
	}
}

// handleAirspaceNewPost handles requests to create a new airspace with the given form parameters. If the
// request is successful, the user will be redirected to the airspace listings page.
func handleAirspaceNewPost(c *gin.Context) {
	client := getCaptainClient()
	// TODO: Validate input.
	_, err := client.CreateAirspace(c.PostForm("HumanName"), c.PostForm("NetName"))
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
	} else {
		c.Redirect(http.StatusFound, "/airspace")
	}
}

// handleAirspaceDelete handles requests to delete an airspace of the given ID. If the request is
// successful, the user will be redirected to the airspace listings page.
func handleAirspaceDelete(c *gin.Context) {
	client := getCaptainClient()
	airspace, err := getAirspaceFromURLParameter(c, client)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	err = client.DeleteAirspace(airspace.ID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return
	}
	c.Redirect(http.StatusFound, "/airspace")
}
