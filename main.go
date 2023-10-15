package main

import (
	"context"
	"fmt"
	"github.com/dark-enstein/port/config"
	"github.com/dark-enstein/port/db"
	"github.com/dark-enstein/port/server"
	"github.com/dark-enstein/port/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	S             = server.NewService()
	ConfigInvalid = fmt.Errorf("passed in config is invalid. check logs")
)

func InitSys() chan os.Signal {
	S.Ctx = context.Background()
	S.Cfg = config.NewConfig()
	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM)
	return cancel
}

func SetStage() error {

	//config.ResolveConfig()
	S.Cfg = config.NewConfig()
	set := config.CommandLine
	set.StringVar(&S.Cfg.LogLevel, config.FlagLogLevel, config.DefaultFLagLogLevel, "-")
	set.StringVar(&S.Cfg.Port, config.FlagPort, config.DefaultFlagPort, "-")
	set.StringVar(&S.Cfg.EnabledDB, config.FlagDB, config.DefaultFlagDB, "-")
	set.StringVar(&S.Cfg.DBHost, config.FlagDBHost, config.DefaultFlagDBHost, "-")
	err := set.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("unable to parse arguments: %w", err)
	}

	S.Log = config.NewLogger(S.Cfg.LogLevel)
	logger := S.Log.With().Str("method", "SetStage()").Logger()
	S.Ctx = context.WithValue(S.Ctx, util.LoggerInContext, S.Log)
	S.Ctx = context.WithValue(S.Ctx, util.ConfigInContext, S.Cfg)

	//validate config
	if !S.ValidateConfig() {
		return ConfigInvalid
	}

	S.DB, err = db.NewClient(S.Ctx, S.Cfg.EnabledDB, S.Cfg.DBHost)
	if err != nil {
		return err
	}

	isConnected := S.DB.Ping()
	if !isConnected {
		logger.Info().Msg("cannot ping db")
	} else {
		logger.Info().Msg("established ping to db")
	}

	return nil
}

func main() {
	cancel := InitSys()
	err := SetStage()
	if err != nil {
		_ = fmt.Errorf("couldn't start server: %w", err)
		os.Exit(1)
	}

	go server.Run() // start server
	if e := S.Log.Debug(); e.Enabled() {
		e.Msgf("registered port %v", S.Cfg.ConstructPort())
	}

	<-cancel
	S.Log.Info().Msg("server stopped")

	// impl graceful shutdown. upload sorta dumps?
	ctx, lastCancel := context.WithTimeout(S.Ctx, 5*time.Second)
	defer func() {
		// extra handling steps; close database, redis, truncate message queues, etc.
		lastCancel()
	}()

	if err := S.Srv.Shutdown(ctx); err != nil {
		S.Log.Info().Msgf("server shutdown failed %v", err)
	}

	S.Log.Info().Msg("server shutdown properly")
}
