package execsrv

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

type heartbeatInformer struct {
	execID    string
	observers map[string]HeartbeatObserver
	mu        sync.Mutex
	driver    Driver
	maxFail   int
	failCount int
}

// NewPollingHeartbeatInformerFunc returns a HeartbeatInformerFunc that polls the driver
// for the exec status. If the driver fails to return the status for maxFail times, the
// informer will inform all observers that the exec has failed and exit.
func NewPollingHeartbeatInformerFunc(driver Driver, maxFail int) HeartbeatInformerFunc {
	// Maintain a map of execs and only create a new informer once per exec
	execs := map[string]*heartbeatInformer{}

	return func(exec types.Exec) HeartbeatInformer {

		if _, ok := execs[exec.ID]; ok {
			return execs[exec.ID]
		}

		inf := &heartbeatInformer{
			execID:    exec.ID,
			observers: make(map[string]HeartbeatObserver),
			mu:        sync.Mutex{},
			driver:    driver,
			maxFail:   maxFail,
			failCount: 0,
		}

		execs[exec.ID] = inf

		return inf
	}
}

func (b *heartbeatInformer) Inform(id string, heartbeat Heartbeat) {
	b.mu.Lock()
	defer b.mu.Unlock()

	o, ok := b.observers[id]
	if !ok {
		return
	}

	go o.Update(heartbeat)
}

func (b *heartbeatInformer) Register(o HeartbeatObserver) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.observers[o.ID()]; ok {
		return
	}
	b.observers[o.ID()] = o
}

func (b *heartbeatInformer) Unregister(o HeartbeatObserver) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.observers, o.ID())
}

func (b *heartbeatInformer) Watch() {

	log.Info().
		Str(types.ExecIDCtxKey, b.execID).
		Msgf("Starting heartbeat informer for exec %s", b.execID)

	go func() {
		defer func() {
			log.Info().
				Str(types.ExecIDCtxKey, b.execID).
				Msgf("Heartbeat informer stopped for exec %s", b.execID)
		}()

		for {
			select {
			case <-time.After(10 * time.Second):

				status, err := b.driver.ExecGetStatus(context.Background(), b.execID)
				if err != nil {
					b.failCount++

					if b.failCount > b.maxFail {
						// No need to notify. Heartbeat just stops.
						log.Info().
							Str(types.ExecIDCtxKey, b.execID).
							Msgf("Heartbeat not detected for exec %q", b.execID)

						return
					}

					continue
				}

				b.failCount = 0

				for _, o := range b.observers {
					b.Inform(o.ID(), Heartbeat{
						ExecID: b.execID,
						Time:   time.Now(),
						Status: status,
					})
				}

			case <-time.After(2 * time.Minute):
				log.Info().
					Str(types.ExecIDCtxKey, b.execID).
					Msgf("Heartbeat informer still running for exec %s", b.execID)
			}
		}
	}()
}

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
