package main

import (
	"github.com/ARMmaster17/Captain/Builder"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {
	log.Debug().Msg("initializing builder service on unified framework")
	builder, err := Builder.NewBuilder()
	if err != nil {
		log.Fatal().Err(err).Stack().Msg("builder did not initialize")
	}
	log.Debug().Msg("starting builder service on new thread")
	builder.StartAsync()
	log.Trace().Msg("setting up interrupt hook")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Debug().Msg("received console interrupt")
	builder.StopAsync()
	log.Trace().Msg("exiting")

}
