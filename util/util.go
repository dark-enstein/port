package util

import (
	"context"
	"fmt"
	"github.com/dark-enstein/port/config"
	"github.com/rs/zerolog"
	"os"
)

const (
	LoggerInContext    = "logger"
	ErrorInContext     = "reqError"
	DBInContext        = "dbConn"
	ConfigInContext    = "serverConfig"
	RequestIDInContext = "requestID"
	QRLocInContext     = "qrLoc"
)

const (
	CREATE = iota
	READ
	UPDATE
	DELETE
	LIST
	UPLOAD
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

// RetrieveReqIDFromCtx returns the request *uuid.UUID stored in the request context
func RetrieveReqIDFromCtx(ctx context.Context) string {
	return ctx.Value(RequestIDInContext).(string)
}

// RetrieveFromCtx is a generic retrieve from context function. It takes in the request context and the key of the value within the contex. The key must be string.
// It returns an any (interface) value which will be type cast by the consumer of the function.
func RetrieveFromCtx(ctx context.Context, key string) any {
	return ctx.Value(RequestIDInContext)
}

// ExitOnErrorln exists with and prints an error to StdErr
// Adapted from original version here: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html
func ExitOnErrorln(msg string) {
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}

// ExitOnErrorf exists and prints an error to StdErr. It takes in standard formatting references.
func ExitOnErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
