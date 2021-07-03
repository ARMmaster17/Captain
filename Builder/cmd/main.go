package main

import (
	"github.com/ARMmaster17/Captain/Builder"
	"github.com/rs/zerolog/log"
)

func main() {
	builder, err := Builder.NewBuilder()
	if err != nil {
		log.Fatal().Err(err).Stack().Msg("builder did not initialize")
	}
	builder.Start()
}
