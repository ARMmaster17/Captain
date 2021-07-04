package framework

import (
	"fmt"
	"github.com/ARMmaster17/Captain/Shared/config"
	"github.com/ARMmaster17/Captain/Shared/logging"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
	log.Trace().Msg("reading environment config")
	config.InitConfiguration(appName)
	logging.InitLogging()
	log.Debug().Str("db", config.GetGlobalString("db")).Msg("connecting to database at CAPTAIN_DB")
	db, err := InitDB(config.GetGlobalString("db"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database:\n%w", err)
	}
	log.Trace().Msg("connected to database")
	return &Framework{
		AppName:   appName,
		Router:    mux.NewRouter(),
		HTTPState: HTTPStopped,
		DB:        db,
	}, nil
}
