package cmd

import (
	"github.com/rs/zerolog"
	log "github.com/thought-machine/falco-probes/internal/logging"
)

// LoggingOpts represents the available logging options for command line tools.
type LoggingOpts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

func configureLoggingFromOpts(opts *LoggingOpts) {
	level := zerolog.InfoLevel
	switch len(opts.Verbose) {
	case 1:
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)
	log.Logger = log.Logger.Level(level)
}
