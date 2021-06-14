package main

import (
	"fmt"
	"github.com/ARMmaster17/Captain/CaptainLib"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

// getCaptainClient is a factory for CaptainClient objects with the base URL injected in. In the future,
// authentication will also be handled here.
func getCaptainClient() *CaptainLib.CaptainClient {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain/radar/")
	err := viper.ReadInConfig()
	if err != nil {
		return CaptainLib.NewCaptainClient("http://localhost:5000/")
	}
	return CaptainLib.NewCaptainClient(viper.GetString("url"))
}

// getURLIDParameter helper method to retrieve an integer value from the request URL.
func getURLIDParameter(name string, c *gin.Context) (int, error) {
	uriName := fmt.Sprintf("%sid", name)
	idVal, err := strconv.ParseInt(c.Param(uriName), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(idVal), nil
}

// forceIntRead helper method to retrieve an integer value from form input. Defaults to zero if an error occurs
// and the error is not passed on to the method caller.
func forceIntRead(input string) int {
	idVal, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0
	}
	return int(idVal)
}

// getAirspaceFromURLParameter returns an Airspace object based on the parameters passed in the request URI.
func getAirspaceFromURLParameter(c *gin.Context, client *CaptainLib.CaptainClient) (CaptainLib.Airspace, error) {
	airspaceID, err := getURLIDParameter("airspace", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid airspace ID:\n%w", err))
		return CaptainLib.Airspace{}, err
	}
	airspace, err := client.GetAirspaceByID(airspaceID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return CaptainLib.Airspace{}, err
	}
	return airspace, nil
}

// getFlightFromURLParameter returns an Airspace object based on the parameters passed in the request URI.
func getFlightFromURLParameter(c *gin.Context, client *CaptainLib.CaptainClient) (CaptainLib.Flight, error) {
	flightID, err := getURLIDParameter("flight", c)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid flight ID:\n%w", err))
		return CaptainLib.Flight{}, err
	}
	flight, err := client.GetFlightByID(flightID)
	if err != nil {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Error:\n%w", err))
		return CaptainLib.Flight{}, err
	}
	return flight, nil
}
