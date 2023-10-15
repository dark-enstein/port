package model

import (
	"context"
	"github.com/dark-enstein/port/util"
	"github.com/rs/zerolog"
)

type Logger zerolog.Logger

// RetrieveLoggerFromCtx returns the *zerolog.Logger stored in the request context
func RetrieveLoggerFromCtx(ctx context.Context) *Logger {
	log := Logger(*ctx.Value(util.LoggerInContext).(*zerolog.Logger))
	return &log
}

func (l *Logger) WithMethod(meth string) *zerolog.Logger {
	zL := (*zerolog.Logger)(l).With().Str("method", meth).Logger()
	return &zL
}
