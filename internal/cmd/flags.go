package cmd

import (
	"os"
	"path"

	"github.com/jessevdk/go-flags"
	"github.com/thought-machine/falco-probes/internal/logging"
)

// MustParseFlags parses the given application options from command line arguments.
func MustParseFlags(opts interface{}) {
	appName := path.Base(os.Args[0])
	flagParser := flags.NewNamedParser(appName, flags.Default)

	flagParser.EnvNamespaceDelimiter = "_"
	flagParser.NamespaceDelimiter = "_"

	loggingOpts := &LoggingOpts{}
	flagParser.AddGroup("logging options", "logging options", loggingOpts)

	flagParser.AddGroup(appName+" options", "", opts)

	args, err := flagParser.Parse()

	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			handleFlagsErr(flagsErr)
		}
		logging.Logger.Fatal().Err(err).Msg("could not parse flags")
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
