package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/config"
	"github.com/dark-enstein/port/db/model"
	"github.com/dark-enstein/port/util"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	LocalMongoHost      = "mongodb://localhost:27017/"
	MongoHost           = ""
	Mongo               = "mongo"
	err                 = errors.New("")
	DefaultDatabaseOpts = options.DatabaseOptions{
		ReadConcern:    nil,
		WriteConcern:   nil,
		ReadPreference: nil,
		BSONOptions:    nil,
		Registry:       nil,
	}
)

type MongoClient struct {
	config struct {
		kind string
		host string
	}
	conn        *mongo.Client
	collections map[string]*mongo.Collection
	databases   map[string]*mongo.Database
	Opts        Opts
	isConnected bool
	ctx         context.Context
	log         *zerolog.Logger
}

type MongoOptions struct {
	Database         string
	Collection       string
	Table            string
	CreateOnNotExist bool
}

func (mo *MongoOptions) IsValid() bool {
	// TODO: impl later
	return true
}

func (m *MongoClient) ClientOptsIsValid() (*MongoClient, error) {
	if !m.Opts.ClientOpts.IsValid() {
		m := new(MongoClient) // clean parent struct on error
		return m, fmt.Errorf("client opts isn't valid")
	}
	return m, nil
}

type Opts struct {
	ServerOpts *options.ClientOptions // server db options
	ClientOpts *MongoOptions          // per request db options
}

func NewOpts() *Opts {
	return &Opts{}
}

func (op *Opts) WithServerOpt(srv *options.ClientOptions) *Opts {
	op.ServerOpts = srv
	return op
}

func (op *Opts) WithClientOpt(cli *MongoOptions) *Opts {
	op.ClientOpts = cli
	return op
}

func NewMongoClient(ctx context.Context, host string) (*MongoClient, error) {
	if host == "" {
		host = LocalMongoHost
	}
	cli := &MongoClient{}
	cli.config.kind = Mongo
	cli.config.host = host
	inited, err := cli.Init(ctx)
	return inited, err
}

// Init initializes mongo client. Takes in the server context
func (m *MongoClient) Init(ctx context.Context) (*MongoClient, error) {
	m.ctx = ctx
	m.log = config.NewLoggerWithDebug()
	m.Opts.ServerOpts = options.Client().ApplyURI(m.config.host)
	m.conn, err = mongo.Connect(m.ctx, m.Opts.ServerOpts)
	if err != nil {
		m.isConnected = false
		log.Fatalln(err)
		return m, nil
	}
	m.isConnected = true

	return m, err
}

func (m *MongoClient) Ping() bool {
	err := m.conn.Ping(m.ctx, nil)
	if err != nil {
		return false
	}
	return true
}

func (m *MongoClient) Kind() string {
	return m.config.kind
}

func (m *MongoClient) Host() string {
	return m.config.host
}

// EnsureDBScaffold returns an error if the target database and collections to be used exist. If overide is true, it creates the scaffold.
func (m *MongoClient) EnsureDBScaffold(ctx context.Context, override bool) error {
	log := RetrieveLoggerFromCtx(ctx, "EnsureDBScaffold()")
	m.collections = make(map[string]*mongo.Collection)
	onlyNames := true
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
		return err
	}
	if len(collSlice) < 1 {
		log.Error().Err(fmt.Errorf("not found: no collections in specified database: %v", m.Opts.ClientOpts.Database))
		if !override {
			return fmt.Errorf("not found: no collections in specified database: %v", m.Opts.ClientOpts.Database)
		}
		m.collections[m.Opts.ClientOpts.Collection] = dbase.Collection(m.Opts.ClientOpts.Collection)
	} else {
		if !util.IsIn(m.Opts.ClientOpts.Collection, collSlice) {
			log.Info().Msgf("not found: collection isn't present in specified database: %v", m.Opts.ClientOpts.Database)
			if !override {
				return fmt.Errorf("not found: collection isn't present in specified database: %v", m.Opts.ClientOpts.Database)
			}
			m.collections[m.Opts.ClientOpts.Collection] = dbase.Collection(m.Opts.ClientOpts.Collection)
		}
	}
	log.Info().Msgf("returned collection %v", m.collections[m.Opts.ClientOpts.Collection].Name())
	return nil
}
