package exec

import (
	"time"

	"github.com/unweave/unweave/api/types"
)

// StateInformer informs observers of state changes in registered execs. There should only
// ever be one StateInformer per driver guaranteeing that exec state change is only ever
// transmitted once.
type StateInformer interface {
	Inform(id string, status types.Status)
	Register(o StateObserver)
	Unregister(o StateObserver)
	Watch()
}

// StateObserver listens for exec state changes and handles them based on the implementation
// of the Update method
type StateObserver interface {
	ID() string
	ExecID() string
	Update(status types.Status)
}

type StateObserverFunc func(exec types.Exec) StateObserver

// Stats represents the resource usage of an exec.
type Stats struct {
	CPU  float64
	Mem  float64
	Disk float64
	GPU  float64
}

// StatsInformer regularly Inform observers of the resource utilization of registered
// execs.
type StatsInformer interface {
	Inform(id string, stats Stats)
	Register(o StatsObserver)
	Unregister(o StatsObserver)
	Watch()
}

// StatsObserver listens for exec stats and updates the exec based on the implementing
// policy
type StatsObserver interface {
	ID() string
	Update(stats Stats)
}

type StatsObserverFunc func(exec types.Exec) StatsObserver

// A Heartbeat represents a signal from an exec indicating its status.
type Heartbeat struct {
	ExecID string
	Time   time.Time
	Status types.Status
}

// HeartbeatInformer informs observers of heartbeats in registered execs.
type HeartbeatInformer interface {
	Inform(id string, heartbeat Heartbeat)
	Register(o HeartbeatObserver)
	Unregister(o HeartbeatObserver)
	Watch()
}

// HeartbeatObserver listens for heartbeats and handles them based on the implementation
// of the Update method.
type HeartbeatObserver interface {
	ID() string
	Update(heartbeat Heartbeat)
}

type HeartbeatObserverFunc func(exec types.Exec) HeartbeatObserver
