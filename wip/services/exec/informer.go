package exec

import (
	"sync"
	"time"

	"github.com/unweave/unweave/api/types"
)

// StateInformer informs observers of state changes in registered execs.
type StateInformer interface {
	Inform(id string, status types.Status)
	Register(o StateObserver)
	Unregister(o StateObserver)
	Watch()
}

type stateInformer struct {
	store     Store
	driver    Driver
	observers map[string]StateObserver
	mu        sync.Mutex
}

func newStateInformer(store Store, driver Driver) *stateInformer {
	return &stateInformer{
		store:     store,
		driver:    driver,
		observers: make(map[string]StateObserver),
	}
}

func (i *stateInformer) Inform(id string, status types.Status) {
	i.mu.Lock()
	defer i.mu.Unlock()

	o := i.observers[id]
	go o.Update(status)
}

func (i *stateInformer) Register(o StateObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.ID()]; ok {
		return
	}
	i.observers[o.ID()] = o
}

func (i *stateInformer) Unregister(o StateObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.ID()]; !ok {
		return
	}
	delete(i.observers, o.ID())
}

func (i *stateInformer) Watch() {
	for {
		select {
		case <-time.After(5 * time.Second):
			for _, o := range i.observers {
				// get state and compare to previous

				exec, err := i.store.Get(o.ID())
				if err != nil {
					// handle error
				}

				// get status of exec from Driver
				// Update with real status
				status := types.StatusRunning
				if status != exec.Status {
					if err = i.store.Update(exec.ID, types.Exec{Status: status}); err != nil {
						// handle error
					}
					// TODO: handle context deadlines
					i.Inform("", status)
				}
			}
		}
	}
}

// StatsInformer regularly Inform observers of the resource utilization of registered
// execs.
type StatsInformer interface {
	Inform(id string, stats Stats)
	Register(o StatsObserver)
	Unregister(o StatsObserver)
	Watch()
}

type statsInformer struct {
	store  Store
	driver Driver
}

func newStatsInformer(store Store, driver Driver) *statsInformer {
	return &statsInformer{
		store:  store,
		driver: driver,
	}
}

func (i *statsInformer) Inform(id string, stats Stats) {

}

func (i *statsInformer) Register(o StatsObserver) {

}

func (i *statsInformer) Unregister(o StatsObserver) {

}

func (i *statsInformer) Watch() {

}
