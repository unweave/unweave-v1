package exec

import (
	"time"

	"github.com/unweave/unweave/api/types"
)

// StateObserver listens for exec state changes and handles them based on the implementation
// of the Update method
type StateObserver interface {
	ID() string
	Update(status types.Status)
}

type stateObserver struct {
	exec types.Exec
	srv  *Service
}

func NewStateObserver(exec types.Exec, srv *Service) StateObserver {
	return &stateObserver{exec: exec, srv: srv}
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

// TerminateIdleObserver watches for exec resource utilization and updates the exec according to
// the configured policy
type TerminateIdleObserver struct {
	exec            types.Exec
	srv             *Service
	lastActive      time.Time
	shouldTerminate bool
	idleTimeout     time.Duration
}

func NewTerminateIdleObserver(exec types.Exec, srv *Service, idleTimeout time.Duration) StatsObserver {
	return &TerminateIdleObserver{
		exec:        exec,
		srv:         srv,
		idleTimeout: idleTimeout,
	}
}

func (o *TerminateIdleObserver) ID() string {
	return o.exec.ID
}

func (o *TerminateIdleObserver) Update(stats Stats) {
	// check utilization is below threshold and idle time is above threshold -> terminate
}
