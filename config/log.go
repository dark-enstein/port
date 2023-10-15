package config

import (
	"github.com/rs/zerolog"
	"os"
)

var (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
	OffLevel   = "off" // TODO implement this later
)

type Log struct {
}

func NewLogger(loglevel string) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	switch loglevel {
	case InfoLevel:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger.Log().Msg("log level is set to info")
	case DebugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger.Log().Msg("log level is set to debug")
	case WarnLevel:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		logger.Log().Msg("log level is set to warn")
	case ErrorLevel:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		logger.Log().Msg("log level is set to error")
	case FatalLevel:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		logger.Log().Msg("log level is set to fatal")
	case PanicLevel:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
		logger.Log().Msg("log level is set to panic")
	}
	return &logger
}

func NewLoggerWithDebug() *zerolog.Logger {
	return NewLogger(DebugLevel)
}

func NewLoggerWithInfo() *zerolog.Logger {
	return NewLogger(InfoLevel)
}

func NewLoggerWithError() *zerolog.Logger {
	return NewLogger(ErrorLevel)
}

func NewLoggerWithWarn() *zerolog.Logger {
	return NewLogger(WarnLevel)
}

func NewLoggerWithFatal() *zerolog.Logger {
	return NewLogger(FatalLevel)
}

func NewLoggerWithPanic() *zerolog.Logger {
	return NewLogger(PanicLevel)
}
