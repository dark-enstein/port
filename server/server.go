package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dark-enstein/port/auth"
	"github.com/dark-enstein/port/config"
	"github.com/dark-enstein/port/db"
	"github.com/dark-enstein/port/internal"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
	"time"
)

var (
	S = &Service{}
)

type Server interface {
	auth.Authentication
	internal.Internal
	//PingDependencies(bool) error // TODO. Implement Ping dependencies
	//ValidateJWT(string) error // TODO. Implement JWT. Validate JWT

	//service.Service // TODO representing all the seperate services in port
	IsLive() bool
	IsReady() bool
}

var (
	StartTime = "ServerStartTime"
)

type Response struct {
	ReqID string `json:"req_id"`
	Time  string `json:"time"`
	Resp  string `json:"response"`
}

func ConstructResponse(reqID, resp string) *Response {
	logger := S.Log.With().Str("method", "ConstructResponse()").Logger()
	r := Response{
		ReqID: reqID,
		Time:  time.Now().String(),
		Resp:  resp,
	}
	logger.Info().Msgf("packaging client response %v", r)
	return &r
}

func (r *Response) MarshalJson() ([]byte, error) {
	return json.Marshal(&r)
}

type Service struct {
	Log *zerolog.Logger
	sync.Mutex
	Ctx   context.Context
	ready bool

	Srv http.Server
	Cfg *config.Config
	r   *mux.Router
	DB  db.DB

	auth.Authentication
	internal.Repository
}

// IsReady checks if the server is currently ready to accept connections
func (s *Service) IsReady() bool {
	return s.ready
}

func (s *Service) IsLive() bool {
	return true
}

// RegisterRoutes registers the servers routes binding it to the server handlers.
func (s *Service) RegisterRoutes() *Service {
	s.r = mux.NewRouter()
	s.r.HandleFunc("/ping", ping).Methods(http.MethodGet)
	s.r.HandleFunc("/register", registerUser).Methods(http.MethodPost)
	s.r.HandleFunc("/generate/{type}", generate).Methods(http.MethodPost)
	//s.r.HandleFunc("/register-tickets", register).Methods(http.MethodPost)
	return s
}

// ping handles calls to the "/ping". It simply responds with "pong".
func ping(resp http.ResponseWriter, req *http.Request) {
	log := S.Log.With().Str("method", "ping()").Logger()
	var b []byte
	log.Info().Msgf("ping body output: %v", req.Body)
	_, _ = req.Body.Read(b)
	log.Info().Msgf("ping byte output: %v", b)
	log.Info().Msgf("ping string(byte) output: %v", string(b))
	_, err := fmt.Fprint(resp, "pong\n")
	if err != nil {
		return
	}
}

func (s *Service) ListenAndServe() {
	log := s.Log.With().Str("method", "ListenAndServe()").Logger()
	s.Srv.Addr = s.Cfg.ConstructPort()
	s.Srv.Handler = s.r
	s.ready = false
	s.Ctx = context.WithValue(s.Ctx, StartTime, time.Now())
	go func() {
		if err := s.Srv.ListenAndServe(); err != nil {
			log.Fatal().Msgf("server startup failed with: (%v)", err)
		}
	}()

	select {
	case <-time.After(5 * time.Second):
		s.ready = true
	}

}

func NewService() *Service {
	return &Service{}
}

// ValidateConfig validates that user config is correct
// it logs an error when one of the configs isn't correct, and returns an appropriate boolean appropriately
func (s *Service) ValidateConfig() bool {
	S = s                                       // reference Service pointer created in main()
	return logLevelIsValid() && dbHostIsValid() // && the rest
}

// Run inits the logger and runs the port service.
func Run() {
	// should print version?

	S.Log.Info().Msgf("starting Port server on %v", S.Cfg.ConstructPort())

	_runServer()
}

func _runServer() {
	S.RegisterRoutes().ListenAndServe()
}
