//nolint:wrapcheck
package volumesrv

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
)

// ServiceRouter is a service that routes requests to the correct provider. In most cases
// use should be using this service instead of provider specific services. This takes care
// of routing requests based on the provider and aggregating responses from multiple
// providers when needed.
type ServiceRouter struct {
	store     Store
	delegates map[types.Provider]Service
}

func NewServiceRouter(store Store, delegates map[types.Provider]Service) Service {
	return &ServiceRouter{
		store:     store,
		delegates: delegates,
	}
}

func (s *ServiceRouter) Provider() types.Provider {
	panic("service router doesn't have a single provider")
}

func (s *ServiceRouter) service(provider types.Provider) Service {
	return s.delegates[provider]
}

func (s *ServiceRouter) Create(ctx context.Context, accountID string, projectID string, provider types.Provider, name string, size int) (types.Volume, error) {
	svc := s.service(provider)
	if svc == nil {
		return types.Volume{}, fmt.Errorf("create: unknown provider: %s", provider)
	}

	return svc.Create(ctx, accountID, projectID, provider, name, size)
}

func (s *ServiceRouter) Delete(ctx context.Context, projectID, idOrName string) error {
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

func (s *ServiceRouter) Get(ctx context.Context, projectID, idOrName string) (types.Volume, error) {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return types.Volume{}, err
	}

	return vol, nil
}

func (s *ServiceRouter) List(ctx context.Context, projectID string) ([]types.Volume, error) {
	vols, err := s.store.VolumeList(projectID)
	if err != nil {
		return nil, err
	}

	return vols, nil
}

func (s *ServiceRouter) Resize(ctx context.Context, projectID, idOrName string, size int) error {
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
