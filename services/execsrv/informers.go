package execsrv

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

type pollingStateInformer struct {
	execID     string
	store      Store
	driver     Driver
	prevStatus types.Status
	observers  map[string]StateObserver
	mu         sync.Mutex
}

func NewPollingStateInformerFunc(store Store, driver Driver) StateInformerFunc {
	return func(exec types.Exec) StateInformer {
		return &pollingStateInformer{
			execID:    exec.ID,
			store:     store,
			driver:    driver,
			observers: make(map[string]StateObserver),
			mu:        sync.Mutex{},
		}
	}

}

func (i *pollingStateInformer) Inform(id string, state State) {
	i.mu.Lock()
	defer i.mu.Unlock()

	o, ok := i.observers[id]
	if !ok {
		return
	}

	log.Info().
		Str(types.ObserverCtxKey, o.Name()).
		Str(types.ExecIDCtxKey, i.execID).
		Msgf("Informing polling state observer of state %q", state.Status)

	go o.Update(state)
}

func (i *pollingStateInformer) Register(o StateObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.ID()]; ok {
		return
	}
	i.observers[o.ID()] = o
}

func (i *pollingStateInformer) Unregister(o StateObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.observers[o.ID()]; !ok {
		return
	}
	delete(i.observers, o.ID())
}

// Watch maintains a cache of the exec state and compares it to both the
// Store and the Driver for changes.
//
// Changes in the Store reflect execs that might have been created or deleted
// since the last watch. Changes in the Driver should only reflect changes in
// the exec's running state i.e. whether it transitioned from initializing to
// running, failed, stopped etc.
func (i *pollingStateInformer) Watch() {
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				// Check the store for changes in the exec's state.
				dbExec, err := i.store.Get(i.execID)
				if err != nil {
					log.Err(err).Msg("failed to get exec from store")
				}

				if dbExec.Status != i.prevStatus {
					i.prevStatus = dbExec.Status
					i.Inform(i.execID, State{Status: dbExec.Status})
				}

			case <-time.After(10 * time.Second):
				// Check the driver for changes in the exec's state.
				status, err := i.driver.ExecGetStatus(context.Background(), i.execID)
				if err != nil {
					log.Err(err).Msg("failed to get exec from driver")
				}

				if status != i.prevStatus {
					i.prevStatus = status
					i.Inform(i.execID, State{Status: status})
				}
			}
		}
	}()
}

type statsInformer struct {
	store  Store
	driver Driver
}

func NewPollingStatsInformerFunc(store Store, driver Driver) StatsInformerFunc {
	return func(exec types.Exec) StatsInformer {
		return &statsInformer{
			store:  store,
			driver: driver,
		}
	}
}

func (i *statsInformer) Inform(id string, stats Stats) {
	log.Ctx(context.Background()).Info().Msg("no-op stats inform method called")
}

func (i *statsInformer) Register(o StatsObserver) {
	log.Ctx(context.Background()).Info().Msg("no-op stats register method called")
}

func (i *statsInformer) Unregister(o StatsObserver) {
	log.Ctx(context.Background()).Info().Msg("no-op stats unregister method called")
}

func (i *statsInformer) Watch() {
	log.Ctx(context.Background()).Info().Msg("no-op stats watch method called")
}
