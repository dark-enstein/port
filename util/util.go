package util

import (
	"context"
	"github.com/dark-enstein/port/config"
	"github.com/rs/zerolog"
)

const (
	LoggerInContext = "logger"
	ErrorInContext  = "reqError"
	DBInContext     = "dbConn"
	ConfigInContext = "serverConfig"
)

var (
	Forbidden = `*@:=/[]?|\"<>+;`
)

// IsIn checks if the bee is in the hive
func IsIn(bee string, hive []string) bool {
	for i := 0; i <= len(hive); i++ {
		if bee == hive[i] {
			return true
		}
	}
	return false
}

// IsInMany checks if any of the bees is in the hive
func IsInMany(bees, hive []string) map[string]bool {
	res := make(map[string]bool, len(bees))
	for i := 0; i < len(bees); i++ {
		res[bees[i]] = IsIn(bees[i], hive)
	}
	return res
}

// FlattenMapToString flattens a map into list of strings
func FlattenMapToString(m map[string]string) []string {
	res := make([]string, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

// Logger holds the logger for port
type Logger zerolog.Logger

// RetrieveLoggerFromCtx returns the *zerolog.Logger stored in the request context
func RetrieveLoggerFromCtx(ctx context.Context) *Logger {
	log := Logger(*ctx.Value(LoggerInContext).(*zerolog.Logger))
	return &log
}

func (l *Logger) WithMethod(meth string) *zerolog.Logger {
	zL := (*zerolog.Logger)(l).With().Str("method", meth).Logger()
	return &zL
}

// RetrieveConfigFromCtx returns the *config.Config stored in the request context
func RetrieveConfigFromCtx(ctx context.Context) *config.Config {
	return ctx.Value(ConfigInContext).(*config.Config)
}
