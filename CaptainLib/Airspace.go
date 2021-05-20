package CaptainLib

import (
	"encoding/json"
	"fmt"
)

// Airspace is an isolated zone of applications and services that run independent of other airspaces.
type Airspace struct {
	ID int
	HumanName string
	NetName string
}

// GetAllAirspaces returns all airpaces managed by the ATC instance.
func (c *CaptainClient) GetAllAirspaces() ([]Airspace, error) {
	results, err := c.restGET("airspaces")
	if err != nil {
		return nil, fmt.Errorf("unable to get list of airspaces:\n%w", err)
	}
	var airspaces []Airspace
	err = json.Unmarshal(results, &airspaces)
	if err != nil {
		return nil, fmt.Errorf("unable to format response as array of Airspaces:\n%w", err)
	}
	return airspaces, nil
}

// GetAirspaceByID returns an Airspace instance from the ATC instance with the given ID (if it exists).
func (c *CaptainClient) GetAirspaceByID(id int) (Airspace, error) {
	results, err := c.restGET(fmt.Sprintf("airspace/%d", id))
	if err != nil {
		return Airspace{}, fmt.Errorf("unable to get list of airspaces:\n%w", err)
	}
	var airspace Airspace
	err = json.Unmarshal(results, &airspace)
	if err != nil {
		return Airspace{}, fmt.Errorf("unable to format response as array of Airspaces:\n%w", err)
	}
	return airspace, nil
}

// CreateAirspace will create an airspace with the given parameters. This airspace will then be managed by the
// connected ATC instance.
func (c *CaptainClient) CreateAirspace(humanName string, netName string) (Airspace, error) {
	result, err := c.restPOST("airspace", map[string]interface{}{
		"HumanName": humanName,
		"NetName": netName,
	})
	if err != nil {
		return Airspace{}, fmt.Errorf("unable to create Airspace:\n%w", err)
	}
	var airspace Airspace
	err = json.Unmarshal(result, &airspace)
	if err != nil {
		return Airspace{}, fmt.Errorf("unable to parse response as an Airspace:\n%w", err)
	}
	return airspace, nil
}

// UpdateAirspace commits an Airspace object to the ATC state database, updating any fields that have changed.
func (c *CaptainClient) UpdateAirspace(id int, humanName string, netName string) error {
	_, err := c.restPUT(fmt.Sprintf("airspace/%d", id), map[string]interface{}{
		"HumanName": humanName,
		"NetName": netName,
	})
	if err != nil {
		return fmt.Errorf("unable to update airspace with ID %d:\n%w", id, err)
	}
	return nil
}

// DeleteAirspace deletes an Airspace, and it will no longer be managed by the connected ATC instance.
func (c *CaptainClient) DeleteAirspace(id int) error {
	resp, err := c.restDELETE(fmt.Sprintf("airspace/%d", id))
	if err != nil {
		return fmt.Errorf("unable to delete airspace with ID %d:\n%w", id, err)
	}
	fmt.Printf("FFFFFFFFFFFFFFFFFFFFF: %s\n", resp)
	return nil
}