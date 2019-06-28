package cmdutil

import (
	"github.com/rs/zerolog/log"
	"os"
)

func Exit(err *error) {
	if *err != nil {
		log.Error().Err(*err).Msg("exited")
		os.Exit(1)
	} else {
		log.Info().Msg("exited")
	}
}
