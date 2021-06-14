package CaptainLib

import (
	"encoding/json"
	"fmt"
)

// Flight is a logical grouping of services that form a single app stack (e.g. database + www + nginx).
type Flight struct {
	ID         int
	AirspaceID int
	Name       string
}

// GetAllFlights returns all flights across all airpsaces that are managed by the connected ATC instance.
func (c *CaptainClient) GetAllFlights() ([]Flight, error) {
	results, err := c.restGET("flights")
	if err != nil {
		return nil, fmt.Errorf("unable to get a list of flights:\n%w", err)
	}
	var flights []Flight
	err = json.Unmarshal(results, &flights)
	if err != nil {
		return nil, fmt.Errorf("unable to format response as array of Flights:\n%w", err)
	}
	return flights, nil
}

// GetFlightsByAirspace returns all flights in the specified airspace.
func (c *CaptainClient) GetFlightsByAirspace(airspaceID int) ([]Flight, error) {
	results, err := c.restGET(fmt.Sprintf("airspace/%d/flights", airspaceID))
	if err != nil {
		return nil, fmt.Errorf("unable to get a list of flights:\n%w", err)
	}
	var flights []Flight
	err = json.Unmarshal(results, &flights)
	if err != nil {
		return nil, fmt.Errorf("unable to format response as array of Flights:\n%w", err)
	}
	return flights, nil
}

// GetFlightByID returns a Flight object that is managed by the connected ATC cluster with the given
// ID (if it exists).
func (c *CaptainClient) GetFlightByID(id int) (Flight, error) {
	results, err := c.restGET(fmt.Sprintf("flight/%d", id))
	if err != nil {
		return Flight{}, fmt.Errorf("unable to get flight by id %d:\n%w", id, err)
	}
	var flight Flight
	err = json.Unmarshal(results, &flight)
	if err != nil {
		return Flight{}, fmt.Errorf("unable to format response as a Flight:\n%w", err)
	}
	return flight, nil
}

// CreateFlight creates a flight that will be managed by the connected ATC instance within the given airsapce.
func (c *CaptainClient) CreateFlight(name string, airspaceID int) (Flight, error) {
	result, err := c.restPOST("flight", map[string]interface{}{
		"AirspaceID": airspaceID,
		"Name":       name,
	})
	if err != nil {
		return Flight{}, fmt.Errorf("unable to create Flight:\n%w", err)
	}
	var flight Flight
	err = json.Unmarshal(result, &flight)
	if err != nil {
		return Flight{}, fmt.Errorf("unable to parse response as a Flight:\n%w", err)
	}
	return flight, nil
}

// UpdateFlight sends the given flight object to the connected ATC instance, and any changes to the Flight
// object will be committed to the state database.
func (c *CaptainClient) UpdateFlight(flight Flight) error {
	_, err := c.restPUT(fmt.Sprintf("flight/%d", flight.ID), map[string]interface{}{
		"AirspaceID": flight.AirspaceID,
		"Name":       flight.Name,
	})
	if err != nil {
		return fmt.Errorf("unable to update flight with ID %d:\n%w", flight.ID, err)
	}
	return nil
}

// DeleteFlight deletes a flight object from the connected ATC instance, and will no longer be managed by ATC.
func (c *CaptainClient) DeleteFlight(id int) error {
	_, err := c.restDELETE(fmt.Sprintf("flight/%d", id))
	if err != nil {
		return fmt.Errorf("unable to delete flight with ID %d:\n%w", id, err)
	}
	return nil
}
