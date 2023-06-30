//nolint:wrapcheck
package providersrv

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type Driver interface {
	ProviderListNodeTypes(ctx context.Context, userID string, filterAvailable bool) ([]types.NodeType, error)
	Provider() types.Provider
}

type ProviderService struct {
	driver Driver
}

func NewProviderService(driver Driver) *ProviderService {
	return &ProviderService{driver: driver}
}

func (s *ProviderService) Provider() types.Provider {
	return s.driver.Provider()
}

func (s *ProviderService) ListNodeTypes(
	ctx context.Context,
	userID string,
	filterAvailable bool,
) ([]types.NodeType, error) {
	return s.driver.ProviderListNodeTypes(ctx, userID, filterAvailable)
}
