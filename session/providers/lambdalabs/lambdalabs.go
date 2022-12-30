package lambdalabs

import (
	"context"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/pkg/random"
	"github.com/unweave/unweave-v2/session/providers/lambdalabs/client"
	"github.com/unweave/unweave-v2/types"
)

const apiURL = "https://cloud.lambdalabs.com/api/v1/"

type InstanceDetails struct {
	Type   client.InstanceType `json:"type"`
	Region client.RegionName   `json:"region"`
	// TODO:
	// 	- Filesystems
}

type Runtime struct {
	InstanceDetails

	client *client.ClientWithResponses
}

func (r *Runtime) AddSSHKey(ctx context.Context, sshKey types.SSHKey) (types.SSHKey, error) {
	if sshKey.PublicKey != nil {
		keys, err := r.ListSSHKeys(ctx)
		if err != nil {
			return types.SSHKey{}, fmt.Errorf("failed to list ssh keys, err: %w", err)
		}
		for _, k := range keys {
			if k.PublicKey != nil && *k.PublicKey == *sshKey.PublicKey {
				log.Info().
					Str(types.RuntimeProviderKey, types.LambdaLabsProvider.String()).
					Msgf("SSH Key %q already exists, using existing key", *sshKey.Name)
				return k, nil
			}
		}
	}
	// Key doesn't exist, create a new one

	if sshKey.Name == nil {
		// This should most like never collide with an existing key, but it is possible.
		// In the future, we should check to see if the key already exists before creating
		// it.
		name := "uw-generated-key-" + random.GenerateRandomPhrase(4, "-")
		sshKey.Name = &name
	}

	log.Info().
		Str(types.RuntimeProviderKey, types.LambdaLabsProvider.String()).
		Msgf("Generating new SSH key %q", *sshKey.Name)

	req := client.AddSSHKeyJSONRequestBody{
		Name:      *sshKey.Name,
		PublicKey: sshKey.PublicKey,
	}
	res, err := r.client.AddSSHKeyWithResponse(ctx, req)
	if err != nil {
		return types.SSHKey{}, err
	}
	if res.JSON200 == nil {
		return types.SSHKey{}, fmt.Errorf("failed to generate SSH key")
	}

	return types.SSHKey{
		Name:      &res.JSON200.Data.Name,
		PublicKey: &res.JSON200.Data.PublicKey,
	}, nil
}

func (r *Runtime) ListSSHKeys(ctx context.Context) ([]types.SSHKey, error) {
	res, err := r.client.ListSSHKeysWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("failed to list ssh keys")
	}

	keys := make([]types.SSHKey, len(res.JSON200.Data))
	for i, k := range res.JSON200.Data {
		keys[i] = types.SSHKey{
			Name:      &k.Name,
			PublicKey: &k.PublicKey,
		}
	}
	return keys, nil
}

func (r *Runtime) InitNode(ctx context.Context, sshKey types.SSHKey) (types.Node, error) {
	// If the SSH key is not provided, generate a new one
	if sshKey.Name == nil {
		k, err := r.AddSSHKey(ctx, sshKey)
		if err != nil {
			return types.Node{}, fmt.Errorf("failed to create a new key, %v", err)
		}
		sshKey = k
	}

	// Launch instance
	log.Info().
		Str(types.RuntimeProviderKey, types.LambdaLabsProvider.String()).
		Msgf("Launching instance with SSH key %q", *sshKey.Name)

	req := client.LaunchInstanceJSONRequestBody{
		FileSystemNames:  nil,
		InstanceTypeName: "",
		Name:             types.Stringy("uw-" + random.GenerateRandomPhrase(3, "-")),
		Quantity:         types.Inty(1),
		RegionName:       "",
		SshKeyNames:      []string{*sshKey.Name},
	}

	res, err := r.client.LaunchInstanceWithResponse(ctx, req)
	if err != nil {
		return types.Node{}, err
	}
	if res.JSON200 == nil {
		return types.Node{}, fmt.Errorf("failed to launch instance")
	}
	if len(res.JSON200.Data.InstanceIds) == 0 {
		return types.Node{}, fmt.Errorf("failed to launch instance")
	}
	return types.Node{
		ID:      res.JSON200.Data.InstanceIds[0],
		KeyPair: sshKey,
		Status:  types.StatusInitializingNode,
	}, nil
}

func (r *Runtime) TerminateNode(ctx context.Context, nodeID string) error {
	return nil
}

func NewProvider(apiKey string) (*Runtime, error) {
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create bearer token provider, err: %v", err)
	}

	llClient, err := client.NewClientWithResponses(apiURL, client.WithRequestEditorFn(bearerTokenProvider.Intercept))
	if err != nil {
		return nil, fmt.Errorf("failed to create client, err: %v", err)
	}

	return &Runtime{
		InstanceDetails: InstanceDetails{},
		client:          llClient,
	}, nil
}
