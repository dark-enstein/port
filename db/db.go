package db

import (
	"context"
	"errors"
	"github.com/dark-enstein/port/db/mongo"
	"github.com/dark-enstein/port/util"
)

var (
	SupportedDBs = []string{Mongo}
	Mongo        = "mongo"
)

func NewClient(ctx context.Context, enabled, host string) (DB, error) {
	dblog := GetLoggerFromCtx(ctx).With().Str("method", "NewClient()").Logger()
	if !util.IsIn(enabled, SupportedDBs) {
		dblog.Info().Msg("enabled client not supported")
		return nil, errors.New("enabled client not supported")
	}
	dblog.Info().Msg("enabled client supported")

	switch enabled {
	case Mongo:
		cli, err := mongo.NewMongoClient(ctx, host)
		return cli, err
	}

	return nil, nil

}
