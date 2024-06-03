package logger

import (
	"os"

	"github.com/rs/zerolog"
)

const KeyReqID = "request_id"

func New(isDebug bool) *zerolog.Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger
}
