package exec

import (
	"github.com/unweave/unweave/api/types"
)

type stateObserver struct {
	exec types.Exec
	srv  *Service
}

func NewStateObserverFunc(s *Service) StateObserverFunc {
	return func(exec types.Exec) StateObserver {
		return &stateObserver{exec: exec, srv: s}
	}
}

func (o *stateObserver) ID() string {
	return o.exec.ID
}

func (o *stateObserver) Update(status types.Status) {
	// the state has been updated
	// act accordingly
	switch status {
	case types.StatusRunning:
		// TODO: Update networking details etc
	}
}
