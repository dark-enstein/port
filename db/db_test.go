package db

import (
	"context"
	"github.com/dark-enstein/port/config"
	"github.com/dark-enstein/port/db/model"
	"github.com/dark-enstein/port/db/mongo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	GlobalTestGeneration = 10
	DefaultTestClient    = "mongo"
)

var (
	ErrorDatabaseUnreachable = `error connecting to %v database at %v`
)

type DBTest struct {
	units TestCreateUnit
	db    DB
	ctx   context.Context
	log   *zerolog.Logger
	suite.Suite
	config struct {
		kind string
		host string
	}
}

// TestCreateUnit defines sample unit list
type TestCreateUnit struct {
	users       []model.User
	targetTable string
}

// InitTestCreateUnit initializes the TestPermission struct
func InitTestCreateUnit() TestCreateUnit {
	var dump TestCreateUnit
	for i := 0; i < GlobalTestGeneration; i++ {
		// gen attributes
	}
	return dump
}

func (s *DBTest) SetupTest() {
	s.log = config.NewLoggerWithDebug()
	s.log.Info().Msg("Starting tests...")
	s.ctx = context.Background()
	var err error
	s.db, err = NewClient(s.ctx, DefaultTestClient, mongo.LocalMongoHost)
	s.config.kind = s.db.Kind()
	s.config.host = s.db.Host()
	s.Assert().NoError(err)

	s.log.Info().Msg("Tests startup complete...")
}

// TestIsIn tests that IsIn function works as expected
func (s *DBTest) TestPermission() {
	log := s.log
	s.units = InitTestCreateUnit()
	log.Debug().Msg("created test table")
	a, b, c, d, e, f := int64(110100), int64(101001), int64(01101), int64(0101), int64(10011), int64(11111)
	if len(s.units.users) == 0 {
		s.units = TestCreateUnit{users: []model.User{
			{
				Name: model.Name{
					FirstName: "ayobami",
					LastName:  "bamigboye",
				},
				Birth: "22/07/1999",
				Roles: []map[string]*int64{
					{
						"gama":    &a,
						"beta":    &b,
						"sinja":   &c,
						"delta":   &d,
						"epsilon": &e,
						"finha":   &f,
					},
				},
			},
		},
			targetTable: "users_test",
		}
	}

	for i := 0; i < len(s.units.users); i++ {
		dbResp := s.db.Create(s.units.users[i], model.CreateOpts{TargetTable: "user"})
		s.Assert().NoError(dbResp.Err)
		s.Assert().NotEmpty(dbResp.ID)
	}
}

// TestPingDB tests the connection integrity to the database by pinging
func (s *DBTest) TestPingDB() {
	log := s.log
	isConnected := s.db.Ping()
	s.Assert().Truef(isConnected, ErrorDatabaseUnreachable, s.config.kind, s.config.host)
	log.Info().Msgf("connection confirmed to %v database: %v", s.config.kind, isConnected)
}

func (s *DBTest) TearDownSuite() {
	log.Info().Msg("Commencing test cleanup")
	//err := cleanUpAfterCatTest()
	//s.Require().NoError(err)
	log.Info().Msg("All testing complete")
}

func TestUtilTest(t *testing.T) {
	suite.Run(t, new(DBTest))
}

//func cleanUpAfterCatTest() error {
//	err := cleanUpAfterTest()
//	// cat content
//	return err
//}
