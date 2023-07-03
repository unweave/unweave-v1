package execsrv

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

type pollingStateInformer struct {
	execID       string
	store        Store
	driver       Driver
	prevStatus   types.Status
	observers    map[string]StateObserver
	mu           sync.Mutex
	manager      *PollingStateInformerManager // for removing itself from the manager
	pollInterval time.Duration
}

type PollingStateInformerManager struct {
	store     Store
	driver    Driver
	informers map[string]*pollingStateInformer

	PollInterval time.Duration
}

// NewPollingStateInformerManager returns a new PollingStateInformerManager that allows for
// adding and removing StateInformer(s) for execs. The informer polls for changes in the
// exec status and notifies all subscribed observers.
func NewPollingStateInformerManager(store Store, driver Driver) *PollingStateInformerManager {
	return &PollingStateInformerManager{
		store:     store,
		driver:    driver,
		informers: make(map[string]*pollingStateInformer),
	}
}

func (m *PollingStateInformerManager) Add(exec types.Exec) StateInformer {
	if _, ok := m.informers[exec.ID]; ok {

		log.Warn().
			Str(types.ExecIDCtxKey, exec.ID).
			Msgf("State informer already exists for exec %s", exec.ID)

		return m.informers[exec.ID]
	}

	inf := &pollingStateInformer{
		execID:    exec.ID,
		store:     m.store,
		driver:    m.driver,
		observers: make(map[string]StateObserver),
		mu:        sync.Mutex{},
		manager:   m,
	}

	m.informers[exec.ID] = inf

	log.Info().
		Str(types.ExecIDCtxKey, exec.ID).
		Msgf("Added state informer for exec %q", exec.ID)

	return inf
}

func (m *PollingStateInformerManager) Remove(execID string) {
	if _, ok := m.informers[execID]; !ok {

		log.Warn().
			Str(types.ExecIDCtxKey, execID).
			Msgf("State informer does not exist for exec %s", execID)

		return
	}
	delete(m.informers, execID)
}

func (i *pollingStateInformer) inform(state State) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, observer := range i.observers {
		log.Info().
			Str(types.ObserverCtxKey, observer.Name()).
			Str(types.ExecIDCtxKey, i.execID).
			Msgf("Informing polling state observer of state %q", state.Status)

		go func(o StateObserver) {
			newState := o.Update(state)
			if newState == state {
				return
			}

			i.inform(newState)
		}(observer)
	}
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
	log.Info().
		Str(types.ExecIDCtxKey, i.execID).
		Msgf("Starting watch for state informer for exec %s", i.execID)

	go func() {
		defer func() {
			log.Info().
				Str(types.ExecIDCtxKey, i.execID).
				Msgf("State informer stopped for exec %s", i.execID)

			i.manager.Remove(i.execID)
		}()

		for {
			select {
			case <-time.After(i.pollInterval):
				// Check the store for changes in the exec's state.
				dbExec, err := i.store.Get(i.execID)
				if err != nil {
					log.Err(err).Msg("failed to get exec from store")
				}

				if dbExec.Status != i.prevStatus {
					log.Info().Msgf("store informing transition %s => %s", i.prevStatus, dbExec.Status)
					i.prevStatus = dbExec.Status

					state := State{Status: dbExec.Status}
					i.inform(state)
				}

				if dbExec.Status.IsTerminal() {
					return
				}

			case <-time.After(i.pollInterval * 2):
				// Check the driver for changes in the exec's state.
				status, err := i.driver.ExecGetStatus(context.Background(), i.execID)
				if err != nil {
					log.Err(err).Msg("failed to get exec from driver")
				}

				if status != i.prevStatus {
					log.Info().Msgf("driver informing transition %s => %s", i.prevStatus, status)
					i.prevStatus = status

					state := State{Status: status}
					i.inform(state)
				}
			}
		}
	}()
}
