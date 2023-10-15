package auth

import (
	"context"
	"github.com/dark-enstein/port/db"
	"github.com/dark-enstein/port/db/model"
	"github.com/dark-enstein/port/util"
	"github.com/rs/zerolog"
)

func GetLoggerFromCtx(ctx context.Context) *zerolog.Logger {
	return ctx.Value(util.LoggerInContext).(*zerolog.Logger)
}

func GetDBFromCtx(ctx context.Context) db.DB {
	return ctx.Value(util.DBInContext).(db.DB)
}

func resolveOpts(kind string) Options {
	switch kind {
	case KindUser:
		return &model.UserOptions{
			Database:         UserDB,
			Collection:       UserCollection,
			Table:            UserTable,
			CreateOnNotExist: true,
		}
	}
	return nil
}
