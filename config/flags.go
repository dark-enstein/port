package config

import (
	"flag"
	"github.com/dark-enstein/port/db/model"
)

const (
	FlagLogLevel   = "log-level"
	FlagPort       = "port"
	FlagDB         = "db"
	FlagDBHost     = "db-host"
	FlagProvider   = "provider"
	FlagLOC        = ""
	NoFlagLogLevel = ""
)

var (
	DefaultFLagLogLevel = "info"
	DefaultFlagDB       = "mongo"
	DefaultFlagDBHost   = "mongodb://localhost:27017/"
	DefaultFlagPort     = "8090"
	DefaultDBName       = model.UserDB
	DefaultFlagProvider = "aws"
	DefaultJWTSecretLOC = "./.jwt/secret"
	DefaultFlagLOC      = "~/.aws/credentials"
)

var (
	LogLevels = []string{"info", "debug", "warn"}
)

var CommandLine = flag.NewFlagSet("port", flag.ExitOnError)
