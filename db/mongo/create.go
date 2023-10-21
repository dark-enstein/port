package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/db/model"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

var (
	ErrUnitDoesntExist            = "inferred unit doesn't exist: "
	ErrCollectionNotFoundInClient = "collection not found in referred MongoClient"
)

// Create creates the unit argument in mongo
func (m *MongoClient) Create(ctx context.Context, unit model.Unit, opts model.Opts) *model.DBResponse { // TODO: eliminate reflection, and properly use interfaces
	llog := RetrieveLoggerFromCtx(ctx, "Create()")
	m.ctx = ctx
	switch unit.Kind() {
	case model.UnitUser:
		u := unit.(*model.User)
		m.Opts.ClientOpts = &MongoOptions{
			Database:         opts.RetrieveDatabase(),
			Collection:       opts.RetrieveCollection(),
			CreateOnNotExist: opts.RetrieveOverride(),
		}
		err := m.EnsureDBScaffold(m.ctx, true)

		// if user exists return an err: already created
		exists, resp := m.IsExist(ctx, unit, m.Opts.ClientOpts.Collection)
		if exists {
			llog.Error().Err(fmt.Errorf("unit to be created already exists in db: %w", resp.Err))
			return resp
		}

		if strings.HasPrefix(resp.Err.Error(), ErrUnitDoesntExist) {
			llog.Error().Err(fmt.Errorf("unit type not recognised: %w", resp.Err))
			return resp
		}

		// generate unique user id
		u.UserID = xid.NewWithTime(time.Now()).String()                             // TODO: Revamp this to pass in the UserID via a function and an argument, not by direct assignment that requires reflection
		one, err := m.collections[m.Opts.ClientOpts.Collection].InsertOne(m.ctx, u) // TODO: Revamp this to use an interface function of CreateFilter to construct the bson.M document of the source unit struct
		if err != nil {
			return &model.DBResponse{Err: err}
		}

		// perform some id cosmetics
		idey := strings.TrimRight(strings.TrimLeft(one.InsertedID.(primitive.ObjectID).String(), `ObjectID(\"`), `\")`)
		llog.Info().Msgf("created record with ID: %s", idey)
		return &model.DBResponse{ID: idey, Err: err}
	}

	llog.Info().Msgf("inferred unit %v doesn't exist", unit.Kind())
	return &model.DBResponse{Err: errors.New("inferred unit doesn't exist")}
}

// IsExist checks if a model.Unit exists within the Collection stored in MongoClient. It returns a boolean and a *model.DBResponse.
// If collection isn't found in MongoClient collectiion slice, it errors out with ErrCollectionNotFoundInClient
func (m *MongoClient) IsExist(ctx context.Context, unit model.Unit, collection string) (bool, *model.DBResponse) {
	llog := RetrieveLoggerFromCtx(ctx, "IsExist()")
	switch unit.Kind() {
	case model.UnitUser:
		u := unit.(*model.User)
		if m.Opts.ClientOpts.Collection == "" {
			m.Opts.ClientOpts = &MongoOptions{
				Collection: collection,
			}
		}

		MongoFindFilter := u.FindFilter().(*primitive.D)
		if _, ok := m.collections[m.Opts.ClientOpts.Collection]; !ok {
			llog.Error().Err(fmt.Errorf("%s: %w", ErrCollectionNotFoundInClient, err))
			return false, &model.DBResponse{Err: errors.New(ErrCollectionNotFoundInClient)}
		}

		find, err := m.collections[m.Opts.ClientOpts.Collection].Find(ctx, MongoFindFilter)
		if err != nil {
			llog.Error().Err(fmt.Errorf("not found: record does not exist in collection: %w", err))
			return false, &model.DBResponse{Err: fmt.Errorf("not found: record does not exist in collection: %w", err)}
		}

		// unmarshal mongo output into bson.M
		var lengthUser bson.M
		err = find.All(ctx, &lengthUser)
		if err != nil {
			llog.Error().Err(fmt.Errorf("could not iterate through found records: %w", err))
			return false, &model.DBResponse{Err: fmt.Errorf("could not iterate through found records: %w", err)}
		}
		if len(lengthUser) > 1 {
			llog.Error().Err(fmt.Errorf("multiple records match the applied filter: %v. records: %#v", MongoFindFilter, lengthUser))
			return false, &model.DBResponse{Err: fmt.Errorf("multiple records match the applied filter: %v. records: %#v", MongoFindFilter, lengthUser)}
		}

		val, ok := lengthUser["user_id"]
		if !ok {
			llog.Error().Err(fmt.Errorf("could not retrieve value from map using indexed key: %s. origin map: %#v", "user_id", lengthUser))
		}

		llog.Info().Msgf("queried record exists with ID: %s", val)
		return true, &model.DBResponse{Err: fmt.Errorf("queried record exists with ID: %s", val)}
	}

	llog.Error().Err(fmt.Errorf(ErrUnitDoesntExist, unit.Kind()))
	return false, &model.DBResponse{Err: fmt.Errorf(ErrUnitDoesntExist, unit.Kind())}
}

func (m *MongoClient) CreateAll(ctx context.Context, unit []model.Unit, opts model.CreateOpts) *model.DBResponse {
	return nil
}
