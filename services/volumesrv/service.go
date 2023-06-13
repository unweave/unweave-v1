package volumesrv

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type Service struct {
	store    Store
	driver   Driver
	provider types.Provider
}

func NewService(store Store, driver Driver) *Service {
	return &Service{
		store:    store,
		driver:   driver,
		provider: driver.VolumeProvider(),
	}
}

func (s *Service) Create(ctx context.Context) (string, error) {
	return "", nil
}

func (s *Service) Delete(ctx context.Context) error {
	return nil
}

func (s *Service) Get(ctx context.Context) (types.Volume, error) {
	return types.Volume{}, nil
}

func (s *Service) List(ctx context.Context) ([]types.Volume, error) {
	return []types.Volume{}, nil
}

func (s *Service) Resize(ctx context.Context) error {
	return nil
}
