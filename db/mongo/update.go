package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

// UpdateResponse is the struct for updating a record from Mongo
type UpdateResponse struct {
	UpsertedCount int64
	MatchedCount  int64
	UpsertedID    string
	ModifiedCount int64
	MapTrain      []bson.M
}

// Update updates the record of model.Unit from mongo database
func (m *MongoClient) Update(ctx context.Context, unit model.Unit, opts model.Opts) *model.DBResponse {
	log := RetrieveLoggerFromCtx(ctx, "Read()")
	m.ctx = ctx
	switch unit.Kind() {
	case model.UnitUser:
		u := unit.(*model.User)
		m.Opts.ClientOpts = &MongoOptions{
			Database:         opts.RetrieveDatabase(),
			Collection:       opts.RetrieveCollection(),
			CreateOnNotExist: opts.RetrieveOverride(),
		}

		err := m.EnsureDBScaffold(ctx, false)
		if err != nil {
			log.Error().Err(fmt.Errorf("error getting collection: %w", err))
			return &model.DBResponse{Err: err}
		}

		collection := m.collections[m.Opts.ClientOpts.Collection]
		if err != nil {
			log.Error().Err(fmt.Errorf("error getting collection: %w", err))
			return &model.DBResponse{Err: err}
		}

		find, err := collection.Find(ctx, u.FindFilter())
		if err != nil {
			log.Error().Err(fmt.Errorf("error finding collection: %w", err))
			return nil
		}
		log.Debug().Msgf("found record: %v", find)

		// bson unmarshal into bson.M slice
		var retrieve UpdateResponse
		retrieve.MapTrain = []bson.M{}
		err = find.All(ctx, retrieve.MapTrain)
		if err != nil {
			log.Error().Err(fmt.Errorf("error reading collection into ReadResponse: %w", err))
			return nil
		}

		// return if update target is more than one
		if len(retrieve.MapTrain) > 1 {
			log.Error().Err(fmt.Errorf("multiple records match the applied filter: %v. records: %#v", u.FindFilter(), retrieve.MapTrain))
			return &model.DBResponse{Err: fmt.Errorf("multiple records match the applied filter: %v. records: %#v", u.FindFilter(), retrieve.MapTrain)}
		}

		one, err := collection.UpdateOne(ctx, u.FindFilter(), u)
		if err != nil {
			return nil
		}
		retrieve.UpsertedCount, retrieve.MatchedCount, retrieve.ModifiedCount, retrieve.UpsertedID = one.UpsertedCount, one.MatchedCount, one.ModifiedCount, strings.TrimRight(strings.TrimLeft(one.UpsertedID.(primitive.ObjectID).String(), `ObjectID(\"`), `\")`)
		log.Debug().Msgf("successfully updated record with id: %v. matched: %v, upserted count: %v, modified: %v", retrieve.UpsertedID, retrieve.MatchedCount, retrieve.UpsertedCount, retrieve.ModifiedCount)
		return &model.DBResponse{Content: retrieve}
	}
	log.Info().Msgf("inferred unit %v doesn't exist", unit.Kind())
	return &model.DBResponse{Err: errors.New("inferred unit doesn't exist")}
}
