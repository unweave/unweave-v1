package exec

import (
	"context"
	"sync"
	"time"

	"github.com/unweave/unweave/api/types"
)

type informer interface {
	inform(ctx context.Context, exec types.Exec)
	register(o observer)
	unregister(o observer)
	watch()
}

type stateInformer struct {
	store     Store
	driver    Driver
	observers map[string]observer
	mu        sync.Mutex
}

func newStateInformer(store Store, driver Driver) *stateInformer {
	return &stateInformer{
		store:     store,
		driver:    driver,
		observers: make(map[string]observer),
	}
}

func (i *stateInformer) inform(ctx context.Context, exec types.Exec) {
	i.mu.Lock()
	o := i.observers[exec.ID]
	i.mu.Unlock()

	o.update(exec)
}

func (i *stateInformer) register(o observer) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.id()]; ok {
		return
	}
	i.observers[o.id()] = o
}

func (i *stateInformer) unregister(o observer) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.id()]; !ok {
		return
	}
	delete(i.observers, o.id())
}

func (i *stateInformer) watch() {
	for {
		select {
		case <-time.After(5 * time.Second):
			for _, o := range i.observers {
				// get state and compare to previous

				exec, err := i.store.Get(o.id())
				if err != nil {
					// handle error
				}

				// get status of exec from Driver
				// update with real status
				status := types.StatusRunning
				if status != exec.Status {
					if err = i.store.Update(exec.ID, types.Exec{Status: status}); err != nil {
						// handle error
					}
					exec.Status = status

					// TODO: handle context deadlines
					go i.inform(context.Background(), exec)
				}
			}
		}
	}
}
