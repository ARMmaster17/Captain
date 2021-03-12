package ampq

type Message struct {
	Operation string `json:"operation"`
	Plane     Plane  `json:"plane"`
	Prep	[]string	`json:"prep"`
}
