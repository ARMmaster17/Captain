package framework

import (
	"fmt"
	"github.com/ARMmaster17/Captain/Shared/config"
	"github.com/ARMmaster17/Captain/Shared/logging"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Framework An object that represents a single running service. Handles all common functions and state storage.
type Framework struct {
	AppName   string
	Router    *mux.Router
	HTTPState HTTPListenStatus
	DB        *gorm.DB
}

// NewFramework Initializes the static environment and returns a new Framework object.
func NewFramework(appName string) (*Framework, error) {
	config.InitConfiguration(appName)
	logging.InitLogging()
	db, err := InitDB(config.GetGlobalString("db"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database:\n%w", err)
	}
	return &Framework{
		AppName:   appName,
		Router:    mux.NewRouter(),
		HTTPState: HTTPStopped,
		DB:        db,
	}, nil
}
