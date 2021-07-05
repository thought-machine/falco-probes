package logging

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is a global logger instance to be used as `logging.Logger.Info()....`` or `var log = logging.Logger; log.Info()....`
var Logger = NewLogger()

// NewLogger returns a new logger instance
func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
