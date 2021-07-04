package logging

import (
	"github.com/ARMmaster17/Captain/Shared/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// InitLogging Initializes the underlying logging platform. Sets up output to stdout or syslog as configured.
func InitLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(config.GetGlobalInt("LOG_LEVEL")))
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})
}
