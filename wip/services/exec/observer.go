package exec

import (
	"github.com/unweave/unweave/api/types"
)

// StateObserver listens for exec state changes and handles them based on the implementation
// of the Update method
type StateObserver interface {
	ID() string
	Update(status types.Status)
}

type StateObserverFunc func(exec types.Exec) StateObserver

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

// StatsObserver listens for exec stats and updates the exec based on the implementing
// policy
type StatsObserver interface {
	ID() string
	Update(stats Stats)
}

type StatsObserverFunc func(exec types.Exec) StatsObserver
