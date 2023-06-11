package providersrv

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type Driver interface {
	ProviderListNodeTypes(ctx context.Context, accountID string, filterAvailable bool) ([]types.NodeType, error)
}

type ProviderService struct {
	driver Driver
}

func NewProviderService(driver Driver) *ProviderService {
	return &ProviderService{driver: driver}
}

func (s *ProviderService) ListNodeTypes(ctx context.Context, accountID string, filterAvailable bool) ([]types.NodeType, error) {
	return s.driver.ProviderListNodeTypes(ctx, accountID, filterAvailable)
}
