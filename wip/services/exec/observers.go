package exec

import (
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

type stateObserver struct {
	exec types.Exec
	srv  *ProviderService
}

func NewStateObserverFunc(s *ProviderService) StateObserverFunc {
	return func(exec types.Exec) StateObserver {
		return &stateObserver{exec: exec, srv: s}
	}
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

func (o *stateObserver) Update(state State) {
	log.Info().Str("exec", o.exec.ID).Msgf("No-op state observer received state update: %s", state.Status)
}
