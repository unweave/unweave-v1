//nolint:wrapcheck
package volumesrv

import (
	"context"
	"fmt"

	"github.com/unweave/unweave-v1/api/types"
)

// DelegatingService is a service that routes requests to the correct provider. In most cases
// use should be using this service instead of provider specific services. This takes care
// of routing requests based on the provider and aggregating responses from multiple
// providers when needed.
type DelegatingService struct {
	store     Store
	delegates map[types.Provider]Service
}

func NewDelegatingService(store Store, services ...Service) Service {
	delegates := make(map[types.Provider]Service)

	for i := range services {
		svc := services[i]
		delegates[svc.Provider()] = svc
	}

	return &DelegatingService{
		store:     store,
		delegates: delegates,
	}
}

func (s *DelegatingService) Provider() types.Provider {
	panic("service router doesn't have a single provider")
}

func (s *DelegatingService) service(provider types.Provider) Service {
	return s.delegates[provider]
}

func (s *DelegatingService) Create(ctx context.Context, accountID string, projectID string, provider types.Provider, name string, size int) (types.Volume, error) {
	svc := s.service(provider)
	if svc == nil {
		return types.Volume{}, fmt.Errorf("create: unknown provider: %s", provider)
	}

	return svc.Create(ctx, accountID, projectID, provider, name, size)
}

func (s *DelegatingService) Delete(ctx context.Context, projectID, idOrName string) error {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return err
	}

	svc := s.service(vol.Provider)
	if svc == nil {
		return fmt.Errorf("get: unknown provider: %s", vol.Provider)
	}

	return svc.Delete(ctx, projectID, idOrName)
}

func (s *DelegatingService) Get(_ context.Context, projectID, idOrName string) (types.Volume, error) {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return types.Volume{}, err
	}

	return vol, nil
}

func (s *DelegatingService) List(_ context.Context, projectID string) ([]types.Volume, error) {
	vols, err := s.store.VolumeList(projectID)
	if err != nil {
		return nil, err
	}

	return vols, nil
}

func (s *DelegatingService) Resize(ctx context.Context, projectID, idOrName string, size int) error {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return err
	}

	svc := s.service(vol.Provider)
	if svc == nil {
		return fmt.Errorf("resize: unknown provider: %s", vol.Provider)
	}

	return svc.Resize(ctx, projectID, idOrName, size)
}
