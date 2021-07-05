package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/thought-machine/falco-probes/internal/logging"
)

// MustParseFlags parses the given application options from command line arguments.
func MustParseFlags(opts interface{}) {
	flagParser := flags.NewNamedParser(filepath.Base(os.Args[0]), flags.Default)

	loggingOpts := &LoggingOpts{}
	flagParser.AddGroup("logging options", "logging options", loggingOpts)
	flagParser.AddGroup("application options", "application options", opts)

	args, err := flagParser.Parse()

	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			handleFlagsErr(flagsErr)
		}
		log.Fatal(err)
	}

	if len(args) > 0 {
		logging.Logger.Fatal().Strs("extra-args", args).Msg("found unexpected extra arguments")
	}

	configureLoggingFromOpts(loggingOpts)
}

func handleFlagsErr(err *flags.Error) {
	if err.Type == flags.ErrHelp {
		os.Exit(0)
	}
}
