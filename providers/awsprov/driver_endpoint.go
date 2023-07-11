package awsprov

import (
	"context"
	"errors"

	"github.com/unweave/unweave/api/types"
)

type EndpointDriver struct{}

func (e *EndpointDriver) EndpointDriverName() string {
	return types.AWSProvider.String()
}

func (e *EndpointDriver) EndpointProvider() types.Provider {
	return types.AWSProvider
}

func (e *EndpointDriver) EndpointCreate(ctx context.Context, project string, endpointID string, execID string, internalPort int32) (string, error) {
	return "", errors.New("endpoints unsupported for aws provider")
}
