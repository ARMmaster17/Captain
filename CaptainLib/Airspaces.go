package CaptainLib

import (
	"encoding/json"
	"fmt"
)

type Airspace struct {
	ID int
	HumanName string
	NetName string
}

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

func (c *CaptainClient) CreateAirspace(humanName string, netName string) (int, error) {
	result, err := c.restPOST("airspace", map[string]interface{}{
		"HumanName": humanName,
		"NetName": netName,
	})
	if err != nil {
		return 0, fmt.Errorf("unable to create Airspace:\n%w", err)
	}
	var airspace Airspace
	err = json.Unmarshal(result, &airspace)
	if err != nil {
		return 0, fmt.Errorf("unable to parse response as an Airspace:\n%w", err)
	}
	return airspace.ID, nil
}

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

func (c *CaptainClient) DeleteAirspace(id int) error {
	_, err := c.restDELETE(fmt.Sprintf("airspace/%d", id))
	if err != nil {
		return fmt.Errorf("unable to delete airspace with ID %d:\n%w", id, err)
	}
	return nil
}