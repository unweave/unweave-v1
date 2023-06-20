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
	"github.com/unweave/unweave/services/execsrv"
	"github.com/unweave/unweave/services/sshkeys"
	"github.com/unweave/unweave/services/volumesrv"
	"github.com/unweave/unweave/tools/gonfig"
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
	execStore := execsrv.NewPostgresStore()

	// TODO: init store
	llDriver, err := lambdalabs.NewAuthenticatedLambdaLabsDriver("")
	if err != nil {
		panic(err)
	}

	llStateInf := execsrv.NewPollingStateInformerManager(execStore, llDriver)
	llStatsInf := execsrv.NewPollingStatsInformerManager(execStore, llDriver)
	llHeartbeatInf := execsrv.NewPollingHeartbeatInformerManager(llDriver, 10)

	volumeStore := volumesrv.NewPostgresStore()
	llVolumeSrv := volumesrv.NewService(volumeStore, llDriver)

	lls := execsrv.NewService(execStore, llDriver, llVolumeSrv, llStateInf, llStatsInf, llHeartbeatInf)
	lls = execsrv.WithStateObserver(lls, execsrv.NewStateObserverFunc(lls))

	if err = lls.Init(); err != nil {
		panic(err)
	}

	execRouter := router.NewExecRouter(execStore, lls, nil)
	sshKeysRouter := router.NewSSHKeysRouter(sshkeys.NewService())

	server.API(cfg, runtimeCfg, execRouter, sshKeysRouter)
}
