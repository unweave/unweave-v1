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
	"github.com/unweave/unweave/providers/awsprov"
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
	volStore := volumesrv.NewPostgresStore()

	lls := lambdaLabsService(execStore, volStore)
	awss := awsService(execStore, volStore)

	delegatingExecSrv := execsrv.NewDelegatingService(execStore, lls, awss)
	execRouter := router.NewExecRouter(execStore, delegatingExecSrv)
	sshKeysRouter := router.NewSSHKeysRouter(sshkeys.NewService())

	server.API(cfg, runtimeCfg, execRouter, sshKeysRouter)
}

func lambdaLabsService(execStore execsrv.Store, volStore volumesrv.Store) execsrv.Service {
	llDriver, err := lambdalabs.NewAuthenticatedLambdaLabsDriver("")
	if err != nil {
		panic(err)
	}

	llStateInf := execsrv.NewPollingStateInformerManager(execStore, llDriver)
	llStatsInf := execsrv.NewPollingStatsInformerManager(execStore, llDriver)
	llHeartbeatInf := execsrv.NewPollingHeartbeatInformerManager(llDriver, 10)

	llVolumeSrv := volumesrv.NewService(volStore, llDriver)

	lls := execsrv.NewService(execStore, llDriver, llVolumeSrv, llStateInf, llStatsInf, llHeartbeatInf)
	lls = execsrv.WithStateObserver(lls, execsrv.NewStateObserverFactory(lls))

	if err = lls.Init(); err != nil {
		panic(err)
	}

	return lls
}

func awsService(execStore execsrv.Store, volStore volumesrv.Store) execsrv.Service {
	ec2, sts, iam, err := awsprov.NewAwsApis("", "", "")
	if err != nil {
		panic(err)
	}

	execDriver := awsprov.NewExecDriverAPI("", "", ec2, sts, iam)
	volDriver := awsprov.NewVolumeDriverAPI("", "", ec2)

	awsStateInf := execsrv.NewPollingStateInformerManager(execStore, execDriver)
	awsStatsInf := execsrv.NewPollingStatsInformerManager(execStore, execDriver)
	awsHeartbeatInf := execsrv.NewPollingHeartbeatInformerManager(execDriver, 10)

	awsVolumeSrv := volumesrv.NewService(volStore, volDriver)

	awss := execsrv.NewService(execStore, execDriver, awsVolumeSrv, awsStateInf, awsStatsInf, awsHeartbeatInf)
	awss = execsrv.WithStateObserver(awss, execsrv.NewStateObserverFactory(awss))

	return awss
}
