package main

import (
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/router"
	"github.com/unweave/unweave/api/server"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/providers/lambdalabs"
	"github.com/unweave/unweave/tools/gonfig"
	execsrv "github.com/unweave/unweave/wip/services/exec"
)

func main() {
	cfg := server.Config{}
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
	db.Q = db.New(conn)

	// Initialize unweave from environment variables
	runtimeCfg := &EnvInitializer{}

	// TODO: init store
	lDriver := lambdalabs.ExecDriver{}
	lStateInf := execsrv.NewPollingStateInformer(nil, lDriver)
	lStatsInf := execsrv.NewPollingStatsInformer(nil, lDriver)

	lls, err := execsrv.NewService(nil, lDriver, lStateInf, lStatsInf)
	if err != nil {
		panic(err)
	}
	lls = execsrv.WithStateObserver(lls, execsrv.NewStateObserverFunc(lls))

	execRouter := router.NewExecRouter(nil, lls, nil)

	server.API(cfg, runtimeCfg, execRouter)
}
