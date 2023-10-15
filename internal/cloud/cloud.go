package cloud

import "context"

type Session interface {
	Kind() string
}

type Client interface {
}

type Cloud interface {
	Do() (string, error)
	BeginSession(context.Context) (*Session, error)
}

type Config struct { // for S3 yet
	provider string
	credLoc  string
	action   struct {
		service string
		verb    string
		target  string
	}
}
