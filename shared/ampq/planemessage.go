package ampq

import "github.com/ARMmaster17/Captain/shared/captain"

type Message struct {
	Operation string `json:"operation"`
	Plane     captain.Plane  `json:"plane"`
	Prep	[]string	`json:"prep"`
}
