package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/auth"
	"github.com/dark-enstein/port/config"
	"github.com/dark-enstein/port/db/mongo"
	"github.com/dark-enstein/port/util"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// dbHostIsValid does the low level validation that the dbHost passed in is valid
// it logs an error if the dbHost config isn't correct
func dbHostIsValid() bool {
	log := S.Log.With().Str("method", "dbHostIsValid()").Logger()
	if S.Cfg.DBHost == mongo.LocalMongoHost || strings.Contains(mongo.LocalMongoHost, S.Cfg.DBHost) {
		return true
	}
	split := strings.Split(S.Cfg.DBHost, ":")
	ip := net.ParseIP(split[0])
	if ip == nil {
		if split[0] != mongo.LocalMongoHost {
			log.Debug().Msgf("host ip passed in: %v isn't valid", split[0])
			return false
		}
	}
	hostInt, err := strconv.Atoi(split[1])
	if err != nil {
		log.Debug().Msgf("host port passed in: %v isn't valid", split[1])
	}
	if len(split) < 2 && hostInt < 1000 && hostInt > 40000 {
		return false
	}
	return true
}

// logLevelIsValid does the low level validation that the loglevel passed in is valid
// it logs an error if the log-level config isn't correct
func logLevelIsValid() bool {
	log := S.Log.With().Str("method", "logLevelIsValid()").Logger()
	if S.Cfg.LogLevel == "" {
		S.Cfg.LogLevel = config.DefaultFLagLogLevel
		return true
	}

	isValid := util.IsIn(S.Cfg.LogLevel, config.LogLevels)

	if !isValid {
		log.Log().Msg("the log level passed in isn't recognized")
		return isValid
	}

	log.Info().Msgf("log level set to %v", S.Cfg.LogLevel)

	return isValid
}

func createUserValidate(j *json.Decoder, resp http.ResponseWriter) (*auth.InternalUser, bool) {
	log := S.Log.With().Str("method", "createUserValidate()").Logger()

	//var b []byte
	//_, err := req.Body.Read(b)
	//log.Info().Msgf("%v", req.Body)
	//if err != nil {
	//	log.Info().Msgf("reading request failed with: %v", err)
	//	return nil, false
	//}

	aga := auth.NewUser()
	log.Info().Msgf("%v", j)
	err := j.Decode(&aga)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		log.Info().Msgf("umarshaling request into json failed with: %v", err)
		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(resp, msg, http.StatusBadRequest)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(resp, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(resp, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(resp, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(resp, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(resp, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Print(err.Error())
			http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return nil, false
	}

	err = j.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		http.Error(resp, msg, http.StatusBadRequest)
		return nil, false
	}

	log.Info().Msgf("User validated: %v %v", aga.Name, aga.Birth)

	return aga.IntoInternal(), !(strings.ContainsAny(aga.IntoInternal().Name().String(), util.Forbidden) && aga.String() == "") // add more validation

}
