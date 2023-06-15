package volumesrv

import (
	"context"
	"fmt"
	"time"

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
	id, err := s.driver.VolumeCreate(ctx, size)
	if err != nil {
		return types.Volume{}, err
	}

	v := types.Volume{
		ID:   id,
		Name: name,
		Size: size,
		State: types.VolumeState{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		Provider: s.provider,
	}

	err = s.store.VolumeAdd(projectID, v.ID, v.Provider)
	if err != nil {
		err = fmt.Errorf("failed to add volume to store: %w", err)

		// Cleanup
		e := s.driver.VolumeDelete(ctx, v.ID)
		if e != nil {
			e = fmt.Errorf("failed to cleanup volume, %w", e)
			err = fmt.Errorf("%s, %w", err, e)
			return types.Volume{}, e

		}

		return types.Volume{}, err
	}

	return v, nil
}

func (s *Service) Delete(ctx context.Context, projectID, idOrName string) error {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return fmt.Errorf("failed to get volume from store: %w", err)
	}

	err = s.driver.VolumeDelete(ctx, vol.ID)
	if err != nil {
		return fmt.Errorf("failed to delete volume: %w", err)
	}

	err = s.store.VolumeDelete(vol.ID)
	if err != nil {
		return fmt.Errorf("failed to delete volume from store: %w", err)
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
		return fmt.Errorf("failed to update volume: %w", err)
	}

	err = s.store.VolumeUpdate(vol.ID, vol)
	if err != nil {
		return fmt.Errorf("failed to update volume in store: %w", err)
	}

	return nil
}