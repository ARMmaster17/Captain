package ampq

type Message struct {
	Operation string `json:"operation"`
	Plane     Plane  `json:"plane"`
	Prep	[]string	`json:"prep"`
}

type Plane struct {
	Name	string	`yaml:"name"`
	CPU		int		`yaml:"cpu"`
	RAM		int		`yaml:"ram"`
	Storage	int		`yaml:"storage"`
}
