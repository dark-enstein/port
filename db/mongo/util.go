package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/db/model"
	"github.com/dark-enstein/port/util"
	"github.com/rs/zerolog"
)

var (
	UserCreateInt = 4 // required entries in table
)

// RetrieveLoggerFromCtx returns the *zerolog.Logger stored in the request context
func RetrieveLoggerFromCtx(ctx context.Context, caller string) *zerolog.Logger {

	log := ctx.Value(util.LoggerInContext).(*zerolog.Logger).With().Str("method", caller).Logger()
	return &log
}

// RetrieveErrorFromCtx returns the error stack stored in the request context
func RetrieveErrorFromCtx(ctx context.Context) error {
	return ctx.Value(util.ErrorInContext).(error)
}

func ErrorWrapInContext(ctx context.Context, err error) context.Context {
	conErr := RetrieveErrorFromCtx(ctx)
	if conErr == nil {
		conErr = errors.New(err.Error())
		return context.WithValue(ctx, util.ErrorInContext, conErr)
	}
	return context.WithValue(ctx, util.ErrorInContext, fmt.Errorf("%s: %w", err.Error(), conErr))
}

func ErrorFromContext(ctx context.Context) error {
	return ctx.Value(util.ErrorInContext).(error)
}

func SetUpUserRepo(u *model.User) *map[string]interface{} {
	var repo = make(map[string]interface{}, UserCreateInt)
	repo["firstname"] = u.Name.FirstName
	repo["lastname"] = u.Name.LastName
	repo["date_of_birth"] = u.Birth
	repo["Roles"] = u.Roles
	return &repo
}
