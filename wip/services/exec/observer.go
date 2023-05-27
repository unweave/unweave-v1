package exec

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type observer interface {
	id() string
	update(ctx context.Context, exec types.Exec)
}

type stateObserver struct {
	exec types.Exec
	srv  *Service
}

func (o *stateObserver) id() string {
	return o.exec.ID
}

func (o *stateObserver) update(ctx context.Context, exec types.Exec) {
	// check if state changed
	if o.exec.Status == exec.Status {
		return
	}
	o.exec = exec

	// act accordingly
	switch exec.Status {
	case types.StatusRunning:
		// TODO: update networking details etc
	}
}
