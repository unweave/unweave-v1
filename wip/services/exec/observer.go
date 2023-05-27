package exec

import (
	"github.com/unweave/unweave/api/types"
)

type observer interface {
	id() string
	update(exec types.Exec)
}

type stateObserver struct {
	exec types.Exec
	srv  *Service
}

func (o *stateObserver) id() string {
	return o.exec.ID
}

func (o *stateObserver) update(exec types.Exec) {
	// the state has been updated
	// act accordingly
	switch exec.Status {
	case types.StatusRunning:
		// TODO: update networking details etc
	}
}
