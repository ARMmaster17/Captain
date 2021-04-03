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
	results, err := DBQueryWithParams("SELECT * FROM airspace WHERE id = ?", strconv.Itoa(id))
	if err != nil {
		return &Airspace{}, fmt.Errorf("unable to find airspace ID %d with error: %w", id, err)
	}
	// Scan results and return first value
}

func (a *Airspace)Save() error {
	database, err := getDBConnection()
	if err != nil {
		return fmt.Errorf("unable to save airspace '%s' with error: %w", a.HumanName, err)
	}
}

func (a *Airspace)Delete() error {
	database, err := getDBConnection()
	if err != nil {
		return fmt.Errorf("unable to delete airspace '%s' with error: %w", a.HumanName, err)
	}
}
