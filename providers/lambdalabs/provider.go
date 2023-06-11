package lambdalabs

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

func (d *Driver) ProviderListNodeTypes(ctx context.Context, accountID string, filterAvailable bool) ([]types.NodeType, error) {
	res, err := d.client.InstanceTypesWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		if res.JSON401 != nil {
			return nil, err401(res.JSON401.Error.Message, nil)
		}
		if res.JSON403 != nil {
			return nil, err403(res.JSON403.Error.Message, nil)
		}
		return nil, errUnknown(res.StatusCode(), nil)
	}

	var nodeTypes []types.NodeType
	for id, data := range res.JSON200.Data {
		data := data

		gpuCount, err := parseGPUCount(id)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse number of GPUs")
			gpuCount = 0
		}

		gpuMem, err := parseGPUMemory(id)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse GPU memory")
			gpuMem = 0
		}

		it := types.NodeType{
			ID:       id,
			Name:     &data.InstanceType.Description,
			Regions:  []string{},
			Price:    &data.InstanceType.PriceCentsPerHour,
			Provider: types.LambdaLabsProvider,
			Specs:    getHardwareSpecFromInstanceTypes(data.InstanceType, gpuMem, gpuCount),
		}

		if filterAvailable && len(data.RegionsWithCapacityAvailable) == 0 {
			continue
		}
		for _, region := range data.RegionsWithCapacityAvailable {
			region := region
			it.Regions = append(it.Regions, region.Name)
		}
		nodeTypes = append(nodeTypes, it)
	}

	return nodeTypes, nil
}
