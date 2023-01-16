package main

import (
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/tools/gonfig"
)

func main() {
	cfg := api.Config{}
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

	store := api.NewMemDB()
	api.API(cfg, runtimeCfg, store)
}
