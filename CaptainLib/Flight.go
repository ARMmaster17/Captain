package CaptainLib

type Flight struct {
	ID         int
	AirspaceID int
	Name       string
}

func (c *CaptainClient) GetAllFlights() ([]Flight, error) {
	results, err := c.restGET("flights")
}


