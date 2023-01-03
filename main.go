package main

import (
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/config"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/pkg/gonfig"
	"github.com/unweave/unweave/runtime"
)

func main() {
	cfg := config.Config{
		APIPort: "8080",
		DB:      config.DBConfig{},
	}
	gonfig.GetFromEnvVariables(&cfg)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	conn, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	dbq := db.New(conn)
	api.API(cfg, &runtime.ConfigFileInitializer{}, dbq)
}
