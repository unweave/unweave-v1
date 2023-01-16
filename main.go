package main

import (
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/runtime"
	server2 "github.com/unweave/unweave/server"
	"github.com/unweave/unweave/tools/gonfig"
)

func main() {
	cfg := server2.Config{}
	gonfig.GetFromEnvVariables(&cfg)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	// Creds store
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get user home directory")
	}

	rcp := filepath.Join(home, ".unweave/runtime-config.json")
	runtimeCfg := &runtime.ConfigFileInitializer{Path: rcp}

	store := server2.NewMemDB()
	server2.Server(cfg, runtimeCfg, store)
}
