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
	llService *VolumeService
	uwService *VolumeService
}

func NewServiceRouter(store Store, lambdalabsService, unweaveService *VolumeService) Service {
	return &ServiceRouter{
		store:     store,
		llService: lambdalabsService,
		uwService: unweaveService,
	}
}

func (s *ServiceRouter) Create(ctx context.Context, accountID string, projectID string, provider types.Provider, name string, size int) (types.Volume, error) {
	switch provider {
	case types.LambdaLabsProvider:
		return s.llService.Create(ctx, accountID, projectID, provider, name, size)
	case types.UnweaveProvider:
		return s.uwService.Create(ctx, accountID, projectID, provider, name, size)
	default:
		return types.Volume{}, fmt.Errorf("unknown provider: %s", provider)
	}
}

func (s *ServiceRouter) Delete(ctx context.Context, projectID, idOrName string) error {
	vol, err := s.store.VolumeGet(projectID, idOrName)
	if err != nil {
		return err
	}

	switch vol.Provider {
	case types.LambdaLabsProvider:
		return s.llService.Delete(ctx, projectID, idOrName)
	case types.UnweaveProvider:
		return s.uwService.Delete(ctx, projectID, idOrName)
	default:
		return fmt.Errorf("unknown provider: %s", vol.Provider)
	}
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

	switch vol.Provider {
	case types.LambdaLabsProvider:
		return s.llService.Resize(ctx, projectID, idOrName, size)
	case types.UnweaveProvider:
		return s.uwService.Resize(ctx, projectID, idOrName, size)
	default:
		return fmt.Errorf("unknown provider: %s", vol.Provider)
	}
}
