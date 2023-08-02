package execsrv

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/db"
)

type stateObserver struct {
	exec types.Exec
	srv  Service
}

func NewStateObserverFactory(s Service) StateObserverFactory {
	fn := func(exec types.Exec) StateObserver {
		return &stateObserver{exec: exec, srv: s}
	}

	return StateObserverFactoryFunc(fn)
}

func (o *stateObserver) ExecID() string {
	return o.exec.ID
}

func (o *stateObserver) ID() string {
	return o.exec.ID
}

func (o *stateObserver) Name() string {
	return "state-observer"
}

func (o *stateObserver) Update(state State) State {
	log.Info().Str("exec", o.exec.ID).Msgf("State observer received state update: %s", state.Status)

	switch state.Status {
	case types.StatusRunning:
		log.Info().
			Str(types.ExecIDCtxKey, o.ExecID()).
			Str(types.ObserverCtxKey, o.Name()).
			Msg("Handling exec state running")

		exec, err := o.srv.RefreshConnectionInfo(context.Background(), o.ExecID())
		if err != nil {
			log.Warn().Err(err).Send()
		}

		o.exec = exec

		update := db.ExecStatusUpdateParams{
			ID:      o.ExecID(),
			Status:  db.UnweaveExecStatus(types.StatusRunning),
			ReadyAt: sql.NullTime{Time: time.Now(), Valid: true},
		}

		if err := db.Q.ExecStatusUpdate(context.Background(), update); err != nil {
			log.Warn().Err(err).Send()
		}

		o.exec.Status = types.StatusRunning
	case types.StatusTerminated:
		log.Info().
			Str(types.ExecIDCtxKey, o.ExecID()).
			Str(types.ObserverCtxKey, o.Name()).
			Msg("Handling exec termination")

		if err := o.srv.Terminate(context.Background(), o.ExecID()); err != nil {
			log.Warn().Err(err).Send()
		}
	case types.StatusPending,
		types.StatusInitializing,
		types.StatusError,
		types.StatusFailed,
		types.StatusSuccess,
		types.StatusUnknown:
	}

	return state
}
