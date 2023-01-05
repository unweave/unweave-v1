package lambdalabs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/pkg/random"
	"github.com/unweave/unweave/providers/lambdalabs/client"
	"github.com/unweave/unweave/types"
)

const apiURL = "https://cloud.lambdalabs.com/api/v1/"

type InstanceDetails struct {
	Type   client.InstanceType `json:"type"`
	Region client.RegionName   `json:"region"`
	// TODO:
	// 	- Filesystems
}

// err400 can happen when ll doesn't have enough capacity to create the instance
func err400(msg string, err error) *types.Error {
	return &types.Error{
		Code:     400,
		Provider: types.LambdaLabsProvider,
		Message:  msg,
		Err:      err,
	}
}

func err401(msg string, err error) *types.Error {
	return &types.Error{
		Code:       401,
		Provider:   types.LambdaLabsProvider,
		Message:    msg,
		Suggestion: "Make sure your LambdaLabs credentials are up to date",
		Err:        err,
	}
}

func err403(msg string, err error) *types.Error {
	return &types.Error{
		Code:       403,
		Provider:   types.LambdaLabsProvider,
		Message:    msg,
		Suggestion: "Make sure your LambdaLabs credentials are up to date",
		Err:        err,
	}
}

func err404(msg string, err error) *types.Error {
	return &types.Error{
		Code:       404,
		Provider:   types.LambdaLabsProvider,
		Message:    msg,
		Suggestion: "",
		Err:        err,
	}
}

func err500(msg string, err error) *types.Error {
	if msg == "" {
		msg = "Unknown error"
	}
	return &types.Error{
		Code:       500,
		Message:    msg,
		Suggestion: "LambdaLabs might be experiencing issues. Check the service status page at https://status.lambdalabs.com/",
		Provider:   types.LambdaLabsProvider,
		Err:        err,
	}
}

// We return this when LambdaLabs doesn't have enough capacity to create the instance.
func err503(msg string, err error) *types.Error {
	return &types.Error{
		Code:     503,
		Provider: types.LambdaLabsProvider,
		Message:  msg,
		Err:      err,
	}
}

func errUnknown(code int, err error) *types.Error {
	return &types.Error{
		Code:       code,
		Message:    "Unknown error",
		Suggestion: "",
		Provider:   types.LambdaLabsProvider,
		Err:        err,
	}
}

type Session struct {
	InstanceDetails

	client *client.ClientWithResponses
}

func (r *Session) GetProvider() types.RuntimeProvider {
	return types.LambdaLabsProvider
}

func (r *Session) AddSSHKey(ctx context.Context, sshKey types.SSHKey) (types.SSHKey, error) {
	if sshKey.PublicKey != nil || sshKey.Name != "" {
		keys, err := r.ListSSHKeys(ctx)
		if err != nil {
			return types.SSHKey{}, fmt.Errorf("failed to list ssh keys, err: %w", err)
		}

		for _, k := range keys {
			if sshKey.Name == k.Name {
				// Key exists, make sure it has the same public key if provided
				if sshKey.PublicKey != nil && *sshKey.PublicKey != *k.PublicKey {
					return types.SSHKey{}, err400("SSH key with the same name already exists with a different public key", nil)
				}
				log.Ctx(ctx).Info().Msgf("SSH Key %q already exists, using existing key", sshKey.Name)
				return k, nil
			}
			if sshKey.PublicKey != nil && *k.PublicKey == *sshKey.PublicKey {
				log.Ctx(ctx).Info().Msgf("SSH Key %q already exists, using existing key", sshKey.Name)
				return k, nil
			}
		}
	}
	// Key doesn't exist, create a new one

	if sshKey.Name == "" {
		// This should most like never collide with an existing key, but it is possible.
		// In the future, we should check to see if the key already exists before creating
		// it.
		name := "uw:" + random.GenerateRandomPhrase(4, "-")
		sshKey.Name = name
	}

	log.Ctx(ctx).Info().Msgf("Generating new SSH key %q", sshKey.Name)

	req := client.AddSSHKeyJSONRequestBody{
		Name:      sshKey.Name,
		PublicKey: sshKey.PublicKey,
	}
	res, err := r.client.AddSSHKeyWithResponse(ctx, req)
	if err != nil {
		return types.SSHKey{}, err
	}
	if res.JSON200 == nil {
		err = fmt.Errorf("failed to generate SSH key")
		if res.JSON401 != nil {
			return types.SSHKey{}, err401(res.JSON401.Error.Message, err)
		}
		if res.JSON403 != nil {
			return types.SSHKey{}, err403(res.JSON403.Error.Message, err)
		}
		if res.JSON400 != nil {
			return types.SSHKey{}, err400(res.JSON400.Error.Message, err)
		}
		return types.SSHKey{}, errUnknown(res.StatusCode(), err)
	}

	return types.SSHKey{
		Name:      res.JSON200.Data.Name,
		PublicKey: &res.JSON200.Data.PublicKey,
	}, nil
}

func (r *Session) ListSSHKeys(ctx context.Context) ([]types.SSHKey, error) {
	log.Ctx(ctx).Info().Msg("Listing SSH keys")

	res, err := r.client.ListSSHKeysWithResponse(ctx)
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

	keys := make([]types.SSHKey, len(res.JSON200.Data))
	for i, k := range res.JSON200.Data {
		k := k
		keys[i] = types.SSHKey{
			Name:      k.Name,
			PublicKey: &k.PublicKey,
		}
	}
	return keys, nil
}

func (r *Session) ListInstanceAvailability(ctx context.Context) ([]types.NodeType, error) {
	log.Ctx(ctx).Info().Msgf("Listing instance availability")

	res, err := r.client.InstanceTypesWithResponse(ctx)
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

	var availableInstanceTypes []types.NodeType
	for id, data := range res.JSON200.Data {
		for _, region := range data.RegionsWithCapacityAvailable {
			it := types.NodeType{
				ID:          id,
				Region:      &region.Name,
				Available:   true,
				Name:        &data.InstanceType.Description,
				Price:       &data.InstanceType.PriceCentsPerHour,
				Description: nil,
				Provider:    types.LambdaLabsProvider,
				Specs: types.NodeSpecs{
					VCPUs:     data.InstanceType.Specs.Vcpus,
					Memory:    data.InstanceType.Specs.MemoryGib,
					GPUMemory: nil,
				},
			}
			availableInstanceTypes = append(availableInstanceTypes, it)
		}
	}
	return availableInstanceTypes, nil
}

func (r *Session) InitNode(ctx context.Context, sshKey types.SSHKey) (types.Node, error) {
	log.Ctx(ctx).Info().Msgf("Launching instance with SSH key %q", sshKey.Name)

	req := client.LaunchInstanceJSONRequestBody{
		FileSystemNames:  nil,
		InstanceTypeName: "gpu_1x_a10",
		Name:             types.Stringy("uw-" + random.GenerateRandomPhrase(3, "-")),
		Quantity:         types.Inty(1),
		RegionName:       "asia-south-1",
		SshKeyNames:      []string{sshKey.Name},
	}

	res, err := r.client.LaunchInstanceWithResponse(ctx, req)
	if err != nil {
		return types.Node{}, err
	}
	if res.JSON200 == nil {
		if res.JSON401 != nil {
			return types.Node{}, err401(res.JSON401.Error.Message, nil)
		}
		if res.JSON403 != nil {
			return types.Node{}, err403(res.JSON403.Error.Message, nil)
		}
		if res.JSON500 != nil {
			return types.Node{}, err500(res.JSON500.Error.Message, nil)
		}
		if res.JSON404 != nil {
			return types.Node{}, err404(res.JSON404.Error.Message, nil)
		}

		// We get a 400 if the instance type is not available. We check for the available
		// instances and return them in the error message. Since this is not critical, we
		// can ignore if there's any errors in the process.
		if res.JSON400 != nil {
			suggestion := ""
			msg := strings.ToLower(res.JSON400.Error.Message)
			if strings.Contains(msg, "not enough capacity") {
				// Get a list of available instances
				instances, e := r.ListInstanceAvailability(ctx)
				if e != nil {
					// Log and continue
					log.Ctx(ctx).Warn().
						Msgf("Failed to get a list of available instances: %v", e)
				}

				b, e := json.MarshalIndent(instances, "", "  ")
				if e != nil {
					log.Ctx(ctx).Warn().
						Msgf("Failed to marshal available instances to JSON: %v", e)
				} else {
					suggestion += string(b)
				}
			}
			e := err503(res.JSON400.Error.Message, nil)
			e.Suggestion = suggestion
			return types.Node{}, e
		}

		return types.Node{}, errUnknown(res.StatusCode(), err)
	}

	if len(res.JSON200.Data.InstanceIds) == 0 {
		return types.Node{}, fmt.Errorf("failed to launch instance")
	}

	return types.Node{
		ID:      res.JSON200.Data.InstanceIds[0],
		KeyPair: sshKey,
		Status:  types.StatusInitializing,
	}, nil
}

func (r *Session) TerminateNode(ctx context.Context, nodeID string) error {
	log.Ctx(ctx).Info().Msgf("Terminating instance %q", nodeID)

	req := client.TerminateInstanceJSONRequestBody{
		InstanceIds: []string{nodeID},
	}
	res, err := r.client.TerminateInstanceWithResponse(ctx, req)
	if err != nil {
		return err
	}
	if res.JSON200 == nil {
		if res.JSON401 != nil {
			return err401(res.JSON401.Error.Message, nil)
		}
		if res.JSON403 != nil {
			return err403(res.JSON403.Error.Message, nil)
		}
		if res.JSON400 != nil {
			return err400(res.JSON400.Error.Message, nil)
		}
		if res.JSON404 != nil {
			return err404(res.JSON404.Error.Message, nil)
		}
		if res.JSON500 != nil {
			return err500(res.JSON500.Error.Message, nil)
		}
		return errUnknown(res.StatusCode(), err)
	}

	return nil
}

func NewSessionProvider(apiKey string) (*Session, error) {
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create bearer token provider, err: %v", err)
	}

	llClient, err := client.NewClientWithResponses(apiURL, client.WithRequestEditorFn(bearerTokenProvider.Intercept))
	if err != nil {
		return nil, fmt.Errorf("failed to create client, err: %v", err)
	}

	return &Session{
		InstanceDetails: InstanceDetails{},
		client:          llClient,
	}, nil
}
