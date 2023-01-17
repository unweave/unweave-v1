package client

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/types"
)

type ProviderService struct {
	client *Client
}

func (p *ProviderService) ListNodeTypes(ctx context.Context, provider types.RuntimeProvider, filterAvailable bool) ([]types.NodeType, error) {
	uri := fmt.Sprintf("providers/%s/node-types", provider)
	query := map[string]string{
		"available": fmt.Sprintf("%t", filterAvailable),
	}
	req, err := p.client.NewAuthorizedRestRequest(Get, uri, query, nil)
	if err != nil {
		return nil, err
	}
	res := &types.NodeTypesListResponse{}
	if err = p.client.ExecuteRest(ctx, req, res); err != nil {
		return nil, err
	}
	return res.NodeTypes, nil
}
