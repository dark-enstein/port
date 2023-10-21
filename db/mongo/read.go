package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/db/model"
	"github.com/dark-enstein/port/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReadResponse is the struct for recording the response from a read operation
type ReadResponse struct {
	MapTrain []bson.M
}

// Read reads the record of model.Unit from mongo database
func (m *MongoClient) Read(ctx context.Context, unit model.Unit, opts model.Opts) *model.DBResponse {
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
		find, err := collection.Find(ctx, u.FindFilter())
		if err != nil {
			log.Error().Err(fmt.Errorf("error finding collection: %w", err))
			return nil
		}
		log.Debug().Msgf("found record: %v", find)

		var retrieve ReadResponse
		retrieve.MapTrain = []bson.M{}
		err = find.All(ctx, retrieve.MapTrain)
		if err != nil {
			log.Error().Err(fmt.Errorf("error reading collection into ReadResponse: %w", err))
			return nil
		}
		log.Debug().Msgf("successfully read found records into Map: %v", retrieve.MapTrain)
		return &model.DBResponse{Content: retrieve}
	}
	log.Info().Msgf("inferred unit %v doesn't exist", unit.Kind())
	return &model.DBResponse{Err: errors.New("inferred unit doesn't exist")}
}

// GetCollection is a helper function to get target mongo.Collection to be worked on
func (m *MongoClient) GetCollection(ctx context.Context) (*mongo.Collection, error) {
	log := RetrieveLoggerFromCtx(ctx, "GetCollection()")

	_ = &options.DatabaseOptions{
		ReadConcern:    nil,
		WriteConcern:   nil,
		ReadPreference: nil,
		BSONOptions:    nil,
		Registry:       nil,
	} // unused for now: TODO

	_ = &options.CollectionOptions{
		ReadConcern:    nil,
		WriteConcern:   nil,
		ReadPreference: nil,
		BSONOptions:    nil,
		Registry:       nil,
	} // unused for now: TODO

	//listDBOpts := &options.ListDatabasesOptions{
	//	NameOnly: &onlyNames,
	//}
	//db, err := m.conn.ListDatabases(ctx, bson.D{{}}, listDBOpts)
	//if err != nil {
	//	return nil, err
	//}

	onlyNames := true
	collectionOpts := &options.ListCollectionsOptions{
		NameOnly: &onlyNames,
	}

	isConnected := m.Ping()
	if !isConnected {
		log.Info().Msg("cannot ping database")
	} else {
		log.Info().Msgf("ping to database successful")
	}

	dbase := m.conn.Database(model.UserDB)

	collSlice, err := dbase.ListCollectionNames(ctx, bson.D{{"options.capped", true}}, collectionOpts)
	log.Info().Msgf("collection slice", collSlice)

	if err != nil {
		log.Error().Err(fmt.Errorf("error listing collection names: %w", err))
		return nil, err
	}

	if len(collSlice) < 1 {
		log.Error().Err(fmt.Errorf("not found: no collections in specified database: %v", m.Opts.ClientOpts.Database))
		return nil, fmt.Errorf("not found: no collections in specified database: %v", m.Opts.ClientOpts.Database)
	}

	if !util.IsIn(m.Opts.ClientOpts.Collection, collSlice) {
		log.Info().Msgf("not found: collection isn't present in specified database: %v", m.Opts.ClientOpts.Database)
	}
	return dbase.Collection(m.Opts.ClientOpts.Collection), nil
}
