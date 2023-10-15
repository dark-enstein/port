package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/dark-enstein/port/internal/cloud"
	"github.com/dark-enstein/port/util"
)

type Config struct {
	credLoc string
	action  struct {
		service string
		verb    string
		target  string
	}
}

type Session struct {
	kind    string
	config  *cloud.Config
	session *session.Session
}

func (s *Session) Kind() string {
	return s.kind
}

var (
	DefaultProfile = "elvis"
)

func (c *Config) NewSessionWithOptions(ctx context.Context) (*Session, error) {
	alog := util.RetrieveLoggerFromCtx(ctx).WithMethod("NewSessionWithOptions()")
	sess, err := session.NewSessionWithOptions(session.Options{Profile: DefaultProfile, Config: aws.Config{
		Region: aws.String("us-west-2"),
	}})
	if err != nil {
		alog.Error().Msgf("creating session with aws failed with %w", err)
		return nil, fmt.Errorf("creating session with aws failed with %w", err)
	}
	alog.Debug().Msg("creating session with aws successful")
	return &Session{kind: "aws", session: sess}, nil
}

func (c *Config) Do() (string, error) {
	return "", nil
}
