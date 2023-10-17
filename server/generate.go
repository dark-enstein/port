package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dark-enstein/port/auth"
	"github.com/dark-enstein/port/internal/generators/qr"
	"github.com/dark-enstein/port/util"
	"github.com/golang/gddo/httputil/header"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	TypeQR = qr.Type("qr")
)

type QR struct {
	Content       string `json:"content"`
	Size          int    `json:"size"`
	RecoveryLevel string `json:"recovery_level"`
}

func NewQR() *QR {
	return &QR{}
}

// generate handles calls to the "/generate". It validates requests and generates a qr code and a link.
func generate(resp http.ResponseWriter, req *http.Request) {
	log := S.Log.With().Str("method", "generate()").Logger()
	ctx := context.WithValue(context.Background(), util.LoggerInContext, S.Log)
	ctx = context.WithValue(ctx, util.DBInContext, S.DB)
	log.Debug().Msg("received a call on /generate, the generate handler is picking it up")

	// if Content-Type header doesn't have its value as "application/json", then return invalid
	// application type
	if req.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(req.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(resp, msg, http.StatusUnsupportedMediaType)
			return
		}
	}
	resp.Header().Set("Content-Type", "application/json")

	// limits the size of the request body to one 1MB.
	//req.Body = http.MaxBytesReader(resp, req.Body, 1048576)
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(time.Second*60))
	defer cancelFunc()

	// validate inputs
	log.Debug().Msg("initiating request validation")
	//isValid := generateValidate(dec, generator, resp)
	qrReq := NewQR()
	err := dec.Decode(&qrReq)
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
	}

	//if err != nil {
	//	resp.WriteHeader(http.StatusBadRequest)
	//	log.Debug().Msg("error with data passed in")
	//	_, err := resp.Write([]byte("data passed in failed validation"))
	//	if err != nil {
	//		log.Debug().Err(fmt.Errorf("error while writing http response %w", err))
	//		_, err := resp.Write([]byte(fmt.Sprint("server error")))
	//		if err != nil {
	//			log.Debug().Err(fmt.Errorf("error while writing http response %w", err))
	//			return
	//		}
	//		return
	//	}
	//	return
	//}
	log.Debug().Msg("successful request validation")

	requestId := uuid.New()
	ctx = context.WithValue(ctx, util.RequestIDInContext, requestId.String())

	var director auth.Director
	switch mux.Vars(req)["type"] {
	case TypeQR.String():
		director = auth.NewQRDirector(ctx, requestId, qrReq.Content, qrReq.RecoveryLevel, qrReq.Size, S.Cfg)
	}

	s, err := director.(*auth.QRDirector).Generate()
	if err != nil {
		log.Error().Msgf("qr generation failed with %v", err)
		genResponse, _ := ConstructErrResponse(requestId.String(), fmt.Sprintf("qr generation failed with %v", err)).MarshalJson()
		_, err := resp.Write(genResponse)
		if err != nil {
			fmt.Fprint(resp, genResponse)
			return
		}
		return
	}
	log.Debug().Msgf("qr generated: %v", s)

	// writing response
	genResponse, err := ConstructResponse(requestId.String(), fmt.Sprintf("generated file at %v", s)).MarshalJson()
	if err != nil {
		log.Error().Msgf("ConstructResponse() failed with %v", err)
		_, err = resp.Write(genResponse)
		return
	}
	resp.WriteHeader(http.StatusOK)
	log.Info().Msgf("file at %v\n", s)
	_, err = resp.Write(genResponse)
	if err != nil {
		fmt.Fprint(resp, genResponse)
		return
	}

	return
}
