package atc

import (
	"fmt"
	"strconv"
)

type Airspace struct {
	AirspaceID int
	HumanName string
	NetName string
}

func CreateAirspace(humanName string, netName string) (*Airspace, error) {
	err := DBExecuteWithParams("INSERT INTO airspace (HumanName, NetName) VALUES (?, ?)", humanName, netName)
	if err != nil {
		return &Airspace{}, fmt.Errorf("unable to create airspace with error: %w", err)
	}
}

func FindByID(id int) (*Airspace, error) {
	results, err := DBQueryWithParams("SELECT * FROM airspace WHERE AirspaceID = ?", strconv.Itoa(id))
	if err != nil {
		return &Airspace{}, fmt.Errorf("unable to find airspace ID %d with error: %w", id, err)
	}
	// Scan results and return first value
	if !results.Next() {
		return &Airspace{}, fmt.Errorf("no airspace found with ID %d", id)
	}
	var idx int
	var humanName string
	var netName string
	err = results.Scan(&idx, &humanName, &netName)
	if err != nil {
		return &Airspace{}, fmt.Errorf("unable to find airspace ID %d with error: %w", id, err)
	}
	return &Airspace{
		AirspaceID: idx,
		HumanName: humanName,
		NetName: netName,
	}, nil
}

func (a *Airspace)Save() error {
	err := DBExecuteWithParams("UPDATE airspace SET HumanName = ?, NetName = ? WHERE AirspaceID = ?", a.HumanName, a.NetName, strconv.Itoa(a.AirspaceID))
	if err != nil {
		return fmt.Errorf("unable to save airspace '%s' with error: %w", a.HumanName, err)
	}
	return nil
}

func (a *Airspace)Delete() error {
	err := DBExecuteWithParams("DELETE FROM airspace WHERE AirspaceID = ?", strconv.Itoa(a.AirspaceID))
	if err != nil {
		return fmt.Errorf("unable to delete airspace '%s' with error: %w", a.HumanName, err)
	}
	return nil
}
