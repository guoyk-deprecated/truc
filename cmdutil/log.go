package cmdutil

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func SetupPlainZerolog(debug bool, color bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out:os.Stdout, NoColor:!color, TimeFormat:time.RFC3339})
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
