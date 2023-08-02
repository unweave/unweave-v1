package volume

import (
	"context"
	"fmt"

	"github.com/unweave/unweave-v1/api/types"
)

type Service struct {
	store     Store
	namespace string
	provider  Provider
}

func NewVolumeService(namespace string, provider Provider, store Store) *Service {
	return &Service{
		store:     store,
		namespace: namespace,
		provider:  provider,
	}
}

func (s *Service) Create(ctx context.Context, size int) (types.Volume, error) {
	vol, err := s.provider.VolumeCreate(ctx, size)
	if err != nil {
		return types.Volume{}, err
	}

	v, err := s.store.Create(s.namespace, vol.ID(), s.provider.Name())
	if err != nil {
		return types.Volume{}, fmt.Errorf("failed to add volume to store: %w", err)
	}

	return v, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.provider.VolumeDelete(ctx); err != nil {
		return fmt.Errorf("failed to delete volume with provider: %w", err)
	}

	if err := s.store.Remove(s.namespace, id); err != nil {
		return fmt.Errorf("failed to remove volume from store: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id string) (types.Volume, error) {
	vol, err := s.store.Get(s.namespace, id)
	if err != nil {
		return types.Volume{}, fmt.Errorf("failed to get volume from store: %w", err)
	}
	return vol, nil
}

func (s *Service) List(ctx context.Context) ([]types.Volume, error) {
	vols, err := s.store.List(s.namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes from store: %w", err)
	}
	return vols, nil
}

func (s *Service) Mount(ctx context.Context, id, path string) error {
	vol, err := s.store.Get(s.namespace, id)
	if err != nil {
		return fmt.Errorf("failed to get volume from store: %w", err)
	}

	v, err := s.provider.VolumeGet(ctx, vol.ID)
	if err := v.Mount(ctx, path); err != nil {
		return fmt.Errorf("failed to mount volume: %w", err)
	}

	return nil
}

func (s *Service) UnMount(ctx context.Context, id, path string) error {
	vol, err := s.store.Get(s.namespace, id)
	if err != nil {
		return fmt.Errorf("failed to get volume from store: %w", err)
	}

	v, err := s.provider.VolumeGet(ctx, vol.ID)
	if err := v.Unmount(ctx, path); err != nil {
		return fmt.Errorf("failed to unmount volume: %w", err)
	}

	return nil
}
