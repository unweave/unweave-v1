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

func (s *Service) Create(ctx context.Context, projectID, name string, size int) (types.Volume, error) {
	err := s.driver.VolumeCreate(ctx, name, size)
	if err != nil {
		return types.Volume{}, err
	}

	v, err := s.driver.VolumeGet(ctx, name)
	if err != nil {
		return types.Volume{}, err
	}

	err = s.store.VolumeAdd(projectID, v)
	if err != nil {
		s.driver.VolumeDelete(ctx, v.ID)
		return types.Volume{}, err
	}

	return v, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.driver.VolumeDelete(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.VolumeDelete(id)
	if err != nil {
		// infers a case where a user is being billed for a volume that does not exist
		return err
	}

	return nil
}

func (s *Service) Get(ctx context.Context, projectID, idOrName string) (types.Volume, error) {
	volume, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return types.Volume{}, err
	}

	return volume, nil
}

func (s *Service) List(ctx context.Context, projectID string) ([]types.Volume, error) {
	vols, err := s.store.VolumeList(projectID)
	if err != nil {
		return vols, err
	}
	return vols, nil
}

func (s *Service) Resize(ctx context.Context, id string, size int) error {
	return nil
}
