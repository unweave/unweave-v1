//nolint:revive
package awsprov

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type ProviderDriver struct {
	supportedInstanceTypes []types.NodeType
}

func NewProviderDriverDefault() *ProviderDriver {
	return NewProviderDriver(nil)
}

func NewProviderDriver(supportedInstanceTypes []types.NodeType) *ProviderDriver {
	return &ProviderDriver{supportedInstanceTypes: supportedInstanceTypes}
}

func (p *ProviderDriver) ProviderListNodeTypes(ctx context.Context, userID string, filterAvailable bool) ([]types.NodeType, error) {
	return p.supportedInstanceTypes, nil
}
