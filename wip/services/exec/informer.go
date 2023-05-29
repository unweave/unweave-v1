package exec

import (
	"fmt"
	"sync"
	"time"

	"github.com/unweave/unweave/api/types"
)

type PollingStateInformer struct {
	store     Store
	driver    Driver
	execs     map[string]types.Exec
	observers map[string]StateObserver
	mu        sync.Mutex
}

func (i *PollingStateInformer) Inform(id string, status types.Status) {
	i.mu.Lock()
	defer i.mu.Unlock()

	o := i.observers[id]
	go o.Update(status)
}

func (i *PollingStateInformer) Register(o StateObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.ID()]; ok {
		return
	}
	i.observers[o.ID()] = o
}

func (i *PollingStateInformer) Unregister(o StateObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.ID()]; !ok {
		return
	}
	delete(i.observers, o.ID())
}

func (i *PollingStateInformer) Watch() {
	for {
		select {
		case <-time.After(5 * time.Second):
			// We need to maintain a cache of execs and compare them to both the Store and
			// the Driver for changes.
			//
			// Changes in the Store reflect execs that might have been created or deleted
			// since the last watch. Changes in the Driver should only reflect changes in
			// the exec's running state i.e. whether it transitioned from initializing to
			// running, failed, stopped etc.
			for _, e := range i.execs {
				// get state and compare to previous

				exec, err := i.store.Get(e.ID)
				if err != nil {
					// handle error
				}

				fmt.Println(exec)

				// if state changed, inform observers

			}
		}
	}
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
