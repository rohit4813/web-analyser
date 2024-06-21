package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// NewLogger returns zerolog.Logger initialized with logger options
func NewLogger(isDebug bool) *zerolog.Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger
}
