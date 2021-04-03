package atc

import "fmt"

type Airspace struct {
	AirspaceID int
	HumanName string
	NetName string
}

func CreateAirspace(humanName string, netName string) (*Airspace, error) {
	database, err := getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to create airspace '%s' with error: %w", humanName, err)
	}
}

func FindByID(id int) (*Airspace, error) {
	database, err := getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to find airspace '%d' with error: %w", id, err)
	}
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
