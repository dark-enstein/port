package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/db/model"
	"strings"
)

type DeletionResponse struct {
	// number of mongo documents deleted
	Count int64
}

func (m *MongoClient) DeleteByID() *model.DBResponse {
	return nil
}

func (m *MongoClient) Delete(ctx context.Context, id string, kind string, unit model.Unit, opts model.Opts) *model.DBResponse {
	llog := RetrieveLoggerFromCtx(ctx, "Create()")
	m.ctx = ctx
	switch kind {
	case model.UnitUser:
		m.Opts.ClientOpts = &MongoOptions{
			Database:         opts.RetrieveDatabase(),
			Collection:       opts.RetrieveCollection(),
			Table:            opts.RetrieveTable(),
			CreateOnNotExist: opts.RetrieveOverride(),
		}
		err := m.EnsureDBScaffold(m.ctx, false)
		if err != nil {
			llog.Error().Err(fmt.Errorf("error checking if id to be deleted exists: %w", err))
			return &model.DBResponse{
				Err: err,
			}
		}

		// check if record exists
		exist, resp := m.IsExist(ctx, unit, m.Opts.ClientOpts.Collection)
		if !exist {
			llog.Error().Err(fmt.Errorf("unit to be deleted doesn't exist in db: %w", resp.Err))
			return resp
		}

		// perform deletion
		one, err := m.collections[m.Opts.ClientOpts.Collection].DeleteOne(ctx, unit.FindFilter())
		if err != nil {
			llog.Error().Err(fmt.Errorf("record deletion errored with: %w", err))
			return &model.DBResponse{Err: err}
		}

		var delResp DeletionResponse
		delResp.Count = one.DeletedCount
		idey := strings.TrimRight(strings.TrimLeft(id, `ObjectID(\"`), `\")`)
		llog.Info().Msgf("created record with ID: %s", idey)
		return &model.DBResponse{ID: idey, Err: err, Content: delResp}
	}

	llog.Info().Msgf("inferred unit %v doesn't exist", unit.Kind())
	return &model.DBResponse{Err: errors.New("inferred unit doesn't exist")}
}
