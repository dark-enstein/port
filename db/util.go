package db

import (
	"context"
	"github.com/dark-enstein/port/util"
	"github.com/rs/zerolog"
)

// RetrieveLoggerFromCtx returns the *zerolog.Logger stored in the request context
func GetLoggerFromCtx(ctx context.Context) *zerolog.Logger {
	return ctx.Value(util.LoggerInContext).(*zerolog.Logger)
}

// RetrieveLoggerFromCtx returns the db.DB stored in the request context
func GetDBFromCtx(ctx context.Context) DB {
	return ctx.Value(util.DBInContext).(DB)
}
