//nolint:revive
package awsprov

import (
	"context"

	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/providers/awsprov/internal/nodes"
)

type ProviderDriver struct {
	supportedInstanceTypes []types.NodeType
}

func NewProviderDriverDefault() *ProviderDriver {
	gNodes := nodes.ToNodeTypesGPU(nodes.GPUNodes())
	cNodes := nodes.CPUNodeTypes()

	return NewProviderDriver(append(gNodes, cNodes))
}

func NewProviderDriver(supportedInstanceTypes []types.NodeType) *ProviderDriver {
	return &ProviderDriver{supportedInstanceTypes: supportedInstanceTypes}
}

func (p *ProviderDriver) Provider() types.Provider {
	return types.AWSProvider
}

func (p *ProviderDriver) ProviderListNodeTypes(ctx context.Context, userID string, filterAvailable bool) ([]types.NodeType, error) {
	return p.supportedInstanceTypes, nil
}
