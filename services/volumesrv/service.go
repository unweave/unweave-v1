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
	volume, err := types.NewVolume(name, size, s.driver.VolumeProvider())
	if err != nil {
		return types.Volume{}, err
	}

	err = s.driver.VolumeCreate(ctx, volume)
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

func (s *Service) Delete(ctx context.Context, projectID, idOrName string) error {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return err
	}

	err = s.driver.VolumeDelete(ctx, vol.ID)
	if err != nil {
		return err
	}

	err = s.store.VolumeDelete(vol.ID)
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

func (s *Service) Resize(ctx context.Context, projectID, idOrName string, size int) error {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return err
	}
	if vol.Size == size {
		return nil
	}

	vol.Size = size

	err = s.driver.VolumeUpdate(ctx, vol)
	if err != nil {
		return nil
	}

	s.store.VolumeUpdate(vol.ID, vol)

	return nil
}
