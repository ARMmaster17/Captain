package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	bootstrapOnly := flag.Bool("boostrap", false, "Runs a stripped-down version of Captain to build the entire Captain stack from a single worker node.")
	apiPort := *flag.Int("apiport", 5000, "Specifies the port to listen on for API requests. Defaults to 5000.")


	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// TODO: Read verbosity level from command line args

	log.Info().Msg("Captain v0.0.0 is starting up")

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
