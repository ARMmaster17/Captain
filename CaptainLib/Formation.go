package CaptainLib

import (
	"encoding/json"
	"fmt"
)

// Formation is a representation of a clonable unit of an application (e.g. database server). Each Formation
// can be scaled up and down to meet demand.
type Formation struct {
	ID          int
	FlightID    int
	Name        string
	CPU         int
	RAM         int
	Disk        int
	BaseName    string
	Domain      string
	TargetCount int
}

// GetAllFormations returns all formations managed by the ATC instance.
func (c *CaptainClient) GetAllFormations() ([]Formation, error) {
	results, err := c.restGET("formations")
	if err != nil {
		return nil, fmt.Errorf("unable to get a list of formations:\n%w", err)
	}
	var formations []Formation
	err = json.Unmarshal(results, &formations)
	if err != nil {
		return nil, fmt.Errorf("unable to format response as array of Formations:\n%w", err)
	}
	return formations, nil
}

// GetFormationByID returns a Formation instance from the conencted ATC instance with the given ID
// (if it exists).
func (c *CaptainClient) GetFormationByID(id int) (Formation, error) {
	results, err := c.restGET(fmt.Sprintf("formation/%d", id))
	if err != nil {
		return Formation{}, fmt.Errorf("unable to get formation with ID %d:\n%w", id, err)
	}
	var formation Formation
	err = json.Unmarshal(results, &formation)
	if err != nil {
		return Formation{}, fmt.Errorf("unable to format response as a Formation:\n%w", err)
	}
	return formation, nil
}

// CreateFormation creates a scalable formation of containers that will be managed by the connected
// ATC instance.
func (c *CaptainClient) CreateFormation(name string, flightID int, CPU int, RAM int, disk int, baseName string, domain string, targetCount int) (Formation, error) {
	result, err := c.restPOST("formation", map[string]interface{}{
		"FlightID": flightID,
		"Name": name,
		"CPU": CPU,
		"RAM": RAM,
		"Disk": disk,
		"BaseName": baseName,
		"Domain": domain,
		"TargetCount": targetCount,
	})
	if err != nil {
		return Formation{}, fmt.Errorf("unable to create Formation:\n%w", err)
	}
	var formation Formation
	err = json.Unmarshal(result, &formation)
	if err != nil {
		return Formation{}, fmt.Errorf("unable to format response as Formation:\n%w", err)
	}
	return formation, nil
}

// UpdateFormation sends the given Formation object to the Captain ATC, and any changes are committed
// to the state database.
func (c *CaptainClient) UpdateFormation(formation Formation) error {
	_, err := c.restPUT(fmt.Sprintf("formation/%d", formation.ID), map[string]interface{}{
		"Name": formation.Name,
		"CPU": formation.CPU,
		"RAM": formation.RAM,
		"Disk": formation.Disk,
		"BaseName": formation.BaseName,
		"Domain": formation.Domain,
		"TargetCount": formation.TargetCount,
	})
	if err != nil {
		return fmt.Errorf("unable to update Formation with ID %d:\n%w", formation.ID, err)
	}
	return nil
}

// DeleteFormation deletes a formation, and it will no longer be managed by the ATC instance.
func (c *CaptainClient) DeleteFormation(id int) error {
	_, err := c.restDELETE(fmt.Sprintf("formation/%d", id))
	if err != nil {
		return fmt.Errorf("unable to delete Formation with ID %d:\n%w", id, err)
	}
	return nil
}
