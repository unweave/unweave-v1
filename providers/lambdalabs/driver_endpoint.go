package lambdalabs

import (
	"context"
	"errors"

	"github.com/unweave/unweave-v1/api/types"
)

type EndpointDriver struct{}

func (e *EndpointDriver) EndpointDriverName() string {
	return types.LambdaLabsProvider.String()
}

func (e *EndpointDriver) EndpointProvider() types.Provider {
	return types.LambdaLabsProvider
}

func (e *EndpointDriver) EndpointCreate(_ context.Context, _, _, _, _ string) (string, error) {
	return "", errors.New("endpoints unsupported for lambdalabs provider")
}
