package lambdalabs

import (
	"context"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/unweave/unweave/providers/lambdalabs/client"
)

// Driver implements the execsrv.Driver and the providersrv.Driver interface for Lambda Labs.
// Lambda Labs needs a special implementation since they don't support Docker on their
// VMs. This means we currently can't run an Exec as a container and instead default to
// the bare VM with the pre-configured Lambda Labs image.
type Driver struct {
	client *client.ClientWithResponses
}

func NewAuthenticatedLambdaLabsDriver(apiKey string) (*Driver, error) {
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create bearer token provider, err: %v", err)
	}

	llClient, err := client.NewClientWithResponses(apiURL, client.WithRequestEditorFn(bearerTokenProvider.Intercept))
	if err != nil {
		return nil, fmt.Errorf("failed to create client, err: %v", err)
	}

	return &Driver{client: llClient}, nil
}

func (d *Driver) ExecPing(ctx context.Context, accountID *string) error {
	_, err := d.client.InstanceTypes(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping LambdaLabs with err %w", err)
	}

	return err
}
