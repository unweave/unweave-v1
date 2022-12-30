package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/api"
	"github.com/unweave/unweave-v2/config"
	"github.com/unweave/unweave-v2/pkg/gonfig"
	"github.com/unweave/unweave-v2/runtime"
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

	api.API(cfg, &runtime.DBInitializer{})
}
