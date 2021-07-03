package framework

import (
	"github.com/ARMmaster17/Captain/Shared/config"
	"github.com/ARMmaster17/Captain/Shared/logging"
	"github.com/gorilla/mux"
)

// Framework An object that represents a single running service. Handles all common functions and state storage.
type Framework struct {
	AppName   string
	Router    *mux.Router
	HTTPState HTTPListenStatus
}

// NewFramework Initializes the static environment and returns a new Framework object.
func NewFramework(appName string) Framework {
	config.InitConfiguration(appName)
	logging.InitLogging()
	return Framework{
		AppName:   appName,
		Router:    mux.NewRouter(),
		HTTPState: HTTPStopped,
	}
}
