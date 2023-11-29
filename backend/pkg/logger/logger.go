package logger

import (
	"os"
	"time"

	stdlog "log"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

var logger zerolog.Logger

// Get returns a pointer to the logger.
func Get() *zerolog.Logger {
	return &logger
}

// InitLogger initializes the logger with given environmentName and logLevelStr
// If the environmentName is not 'debug', log the information in JSON format
func InitLogger(environmentName string, logLevelStr string) error {
	lvl, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(lvl)

	if environmentName == "debug" {
		logger = loggerWithConsoleFormatting()
	} else {
		logger = loggerWithJSONFormatting()
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	zlog.Logger = logger
	return nil
}

func loggerWithConsoleFormatting() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339Nano,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}).With().Timestamp().Caller().Logger()
}

func loggerWithJSONFormatting() zerolog.Logger {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelFieldName = "severity"
	zerolog.MessageFieldName = "message"
	zerolog.CallerFieldName = "caller"
	return zerolog.New(
		os.Stdout,
	).With().Timestamp().Logger()
}
