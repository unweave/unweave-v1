package execsrv

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v1/api/types"
)

type statsInformer struct {
	store  Store
	driver Driver
}

type StatsPollingInformerManager struct {
	store     Store
	driver    Driver
	informers map[string]*statsInformer
}

func NewPollingStatsInformerManager(store Store, driver Driver) *StatsPollingInformerManager {
	return &StatsPollingInformerManager{
		store:     store,
		driver:    driver,
		informers: make(map[string]*statsInformer),
	}
}

func (m *StatsPollingInformerManager) Add(exec types.Exec) StatsInformer {
	if _, ok := m.informers[exec.ID]; ok {
		log.Warn().
			Str(types.ExecIDCtxKey, exec.ID).
			Msgf("Stats informer already exists for exec %s", exec.ID)

		return m.informers[exec.ID]
	}

	inf := &statsInformer{
		store:  m.store,
		driver: m.driver,
	}

	m.informers[exec.ID] = inf

	log.Info().
		Str(types.ExecIDCtxKey, exec.ID).
		Msgf("Stats informer added for exec %s", exec.ID)

	return inf
}

func (m *StatsPollingInformerManager) Remove(execID string) {
	if _, ok := m.informers[execID]; !ok {
		log.Warn().
			Str(types.ExecIDCtxKey, execID).
			Msgf("Stats informer does not exist for exec %s", execID)
		return
	}

	delete(m.informers, execID)
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
