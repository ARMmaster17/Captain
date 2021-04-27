package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	Version string
)

// Main entry point of the application. Handles the creation of the requested number of workers for each task, and sets
// them up to use pipe-based IPC or an external MQ service for communication.
func main() {
	bootstrapOnly := flag.Bool("boostrap", false, "Runs a stripped-down version of Captain to build the entire Captain stack from a single worker node.")
	apiPort := *flag.Int("apiport", 5000, "Specifies the port to listen on for API requests. Defaults to 5000.")


	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// TODO: Read verbosity level from command line args

	log.Info().Msg(fmt.Sprintf("Captain %s is starting up", getApplicationVersion()))

	if *bootstrapOnly {
		err := BootstrapCluster()
		if err != nil {
			log.Error().Msg("An error occurred: " + err.Error())
		}
		return
	}

	apiServer := &APIServer{}
	err := apiServer.Start()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Captain API server has fatally crashed")
		return
	}
	apiServer.Serve(apiPort)

	err = StartMonitoring()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Captain has fatally crashed")
		return
	}
}

// There is a bug in how 'go test' is implemented. This method does not
// have a unit test.
func getApplicationVersion() string {
	return Version
}
