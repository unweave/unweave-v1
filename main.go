package main

import (
	"github.com/rs/zerolog"
	"github.com/unweave/unweave-v2/api"
	"github.com/unweave/unweave-v2/config"
	"github.com/unweave/unweave-v2/pkg/gonfig"
)

func main() {
	cfg := config.Config{
		APIPort: "8080",
		DB:      config.DBConfig{},
	}
	gonfig.GetFromEnvVariables(&cfg)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	api.API(cfg)
}
