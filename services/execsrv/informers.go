package execsrv

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

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
