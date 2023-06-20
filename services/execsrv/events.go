package execsrv

import (
	"time"

	"github.com/unweave/unweave/api/types"
)

type State struct {
	Status types.Status
	Error  error
}

// StateInformer informs observers of state changes in registered execs. There should only
// ever be one StateInformer per driver guaranteeing that exec state change is only ever
// transmitted once.
type StateInformer interface {
	Inform(id string, state State)
	Register(o StateObserver)
	Unregister(o StateObserver)
	Watch()
}

// StateInformerManger manages state informers and ensures that only one informer is
// registered per exec.
type StateInformerManger interface {
	Add(exec types.Exec) StateInformer
	Remove(execID string)
}

// StateObserver listens for exec state changes and handles them based on the implementation
// of the Update method
type StateObserver interface {
	// ID returns the ID of the observer. This should be unique across all observers.
	ID() string
	// ExecID returns the ID of the exec that the observer is observing.
	ExecID() string
	// Name returns the name of the observer and should identify the function of the observer.
	Name() string
	// Update handles the state change of the exec.
	Update(state State)
}

// StateObserverFunc returns a StateObserver for the given exec and StateInformer.
// It takes a reference to the StateInformer to enable passing back state changes to the
// informer.
type StateObserverFunc func(exec types.Exec, inf StateInformer) StateObserver

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

// StatsInformerManger manages stats informers and ensures that only one informer is
// registered per exec.
type StatsInformerManger interface {
	Add(exec types.Exec) StatsInformer
	Remove(execID string)
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

// HeartbeatInformerManger manages heartbeat informers and ensures that only one informer
// is registered per exec.
type HeartbeatInformerManger interface {
	Add(exec types.Exec) HeartbeatInformer
	Remove(execID string)
}

// HeartbeatObserver listens for heartbeats and handles them based on the implementation
// of the Update method.
type HeartbeatObserver interface {
	ID() string
	Update(heartbeat Heartbeat)
}

type HeartbeatObserverFunc func(exec types.Exec) HeartbeatObserver
