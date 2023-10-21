package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dark-enstein/port/auth"
	"github.com/dark-enstein/port/util"
	"github.com/golang/gddo/httputil/header"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// registerUser handles calls to the "/register". It validates requests and creates a user on Port.
func registerUser(resp http.ResponseWriter, req *http.Request) {
	log := S.Log.With().Str("method", "registerUser()").Logger()
	ctx := context.WithValue(context.Background(), util.LoggerInContext, S.Log)
	ctx = context.WithValue(ctx, util.DBInContext, S.DB)
	log.Debug().Msg("received a call on /register, the registerUser handler is picking it up")

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
	req.Body = http.MaxBytesReader(resp, req.Body, 1048576)
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(time.Second*60))
	defer cancelFunc()

	// validate inputs
	log.Debug().Msg("initiating data validation")
	user, isValid := createUserValidate(dec, resp)
	if !isValid {
		resp.WriteHeader(http.StatusBadRequest)
		log.Debug().Msg("error with data passed in")
		_, err := resp.Write([]byte("data passed in failed validation"))
		if err != nil {
			log.Debug().Err(fmt.Errorf("error while writing http response %w", err))
			_, err := resp.Write([]byte(fmt.Sprint("server error")))
			if err != nil {
				log.Debug().Err(fmt.Errorf("error while writing http response %w", err))
				return
			}
			return
		}
		return
	}

	reqId, err := uuid.NewUUID()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		log.Debug().Msgf("generating UUID failed with error: %v", err)
		_, err := resp.Write([]byte(fmt.Sprint("server error")))
		if err != nil {
			log.Debug().Err(fmt.Errorf("error while writing http response %w", err))
			return
		}
	}

	genJWT, err := auth.NewJWT(reqId.String(), "Port", "Register", "Port Inc", time.Duration(5000), S.Cfg.JWTSecretKey)
	if err != nil {
		log.Error().Err(fmt.Errorf("error ehile generating JWT: %w", err))
		return
	}
	resp.Header().Set("Authorization", fmt.Sprintf("Bearer %s", genJWT))
	// store in redis?

	// from here on out copy data into internals
	director := auth.NewUserDirector(ctx)
	log.Debug().Msg("setting up director")

	tx := director.CreateUsers([]auth.InternalUser{*user})
	log.Debug().Msg("tx: ready to execute db call. locking..")
	director.Mutex.Lock()
	ids, err := tx()
	if err != nil {
		log.Debug().Err(fmt.Errorf("error while executing db call: %w", err))
		resp.WriteHeader(http.StatusInternalServerError)
		_, err := resp.Write([]byte(fmt.Sprint("server error")))
		if err != nil {
			log.Debug().Err(fmt.Errorf("error while writing http response %w", err))
			return
		}
		return
	}
	log.Debug().Msg("executed db call. unlocking...")
	director.Mutex.Unlock()

	// build http response
	respBytes, err := ConstructResponse(reqId.String(), fmt.Sprintf("user created with ids: %v", ids), string(genJWT)).MarshalJson()
	log.Debug().Msg("contents of response to client " + string(respBytes))
	if len(respBytes) < 1 || err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		log.Debug().Msg("error with marshalling json into struct")
		_, err := resp.Write([]byte("server error"))
		if err != nil {
			return
		}
		return
	}

	resp.WriteHeader(http.StatusOK)
	_, err = resp.Write(respBytes)
	//_, err = resp.Write([]byte(samp))
	if err != nil {
		log.Debug().Err(fmt.Errorf("failed to marshal api response into response: %w", err))
		resp.WriteHeader(http.StatusInternalServerError)
		_, err := resp.Write([]byte(fmt.Sprint("server error")))
		if err != nil {
			log.Debug().Err(fmt.Errorf("error while writing http response %w", err))
			return
		}
		return
	}

	log.Info().Msgf("user created with id: %s", ids)
	return
}
