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
	manager   *HeartbeatPollingInformerManager // for removing itself from the manager
}

type HeartbeatPollingInformerManager struct {
	driver    Driver
	maxFail   int
	informers map[string]HeartbeatInformer
}

// NewPollingHeartbeatInformerManager returns a new HeartbeatPollingInformerManager that allows for
// adding and removing HeartbeatInformers for execs. The informer polls the driver
// for the exec status and if still active, sends a heartbeat to all subscribed observers.
// If the driver fails to return the status for maxFail times, the informer will exit.
func NewPollingHeartbeatInformerManager(driver Driver, maxFail int) *HeartbeatPollingInformerManager {
	return &HeartbeatPollingInformerManager{
		driver:    driver,
		maxFail:   maxFail,
		informers: make(map[string]HeartbeatInformer),
	}
}

func (h *HeartbeatPollingInformerManager) Add(exec types.Exec) HeartbeatInformer {
	if _, ok := h.informers[exec.ID]; ok {

		log.Warn().
			Str(types.ExecIDCtxKey, exec.ID).
			Msgf("Heartbeat informer already exists for exec %s", exec.ID)

		return h.informers[exec.ID]
	}

	inf := &heartbeatInformer{
		execID:    exec.ID,
		observers: make(map[string]HeartbeatObserver),
		mu:        sync.Mutex{},
		driver:    h.driver,
		maxFail:   h.maxFail,
		failCount: 0,
		manager:   h,
	}

	h.informers[exec.ID] = inf

	log.Info().
		Str(types.ExecIDCtxKey, exec.ID).
		Msgf("Adding heartbeat informer for exec %s", exec.ID)

	return inf
}

func (h *HeartbeatPollingInformerManager) Remove(execID string) {
	if _, ok := h.informers[execID]; !ok {

		log.Warn().
			Str(types.ExecIDCtxKey, execID).
			Msgf("Heartbeat informer does not exist for exec %s", execID)

		return
	}

	delete(h.informers, execID)
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
		Msgf("Starting watch for heartbeat informer for exec %s", b.execID)

	go func() {
		defer func() {
			log.Info().
				Str(types.ExecIDCtxKey, b.execID).
				Msgf("Heartbeat informer stopped for exec %s", b.execID)

			b.manager.Remove(b.execID)
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
