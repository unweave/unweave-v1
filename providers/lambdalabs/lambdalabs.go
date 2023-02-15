package lambdalabs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/providers/lambdalabs/client"
	"github.com/unweave/unweave/tools"
	"github.com/unweave/unweave/tools/random"
)

const apiURL = "https://cloud.lambdalabs.com/api/v1/"

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
	client *client.ClientWithResponses
}

func (s *Session) GetProvider() types.RuntimeProvider {
	return types.LambdaLabsProvider
}

func (s *Session) AddSSHKey(ctx context.Context, sshKey types.SSHKey) (types.SSHKey, error) {
	if sshKey.Name == "" {
		return types.SSHKey{}, fmt.Errorf("SSH key name is required")
	}

	keys, err := s.ListSSHKeys(ctx)
	if err != nil {
		return types.SSHKey{}, fmt.Errorf("failed to list ssh keys, err: %w", err)
	}

	for _, k := range keys {
		if sshKey.Name == k.Name {
			// Key exists, make sure it has the same public key if provided
			if sshKey.PublicKey != nil && *sshKey.PublicKey != *k.PublicKey {
				return types.SSHKey{}, err400("SSH key with the same name already exists with a different public key", nil)
			}
			log.Ctx(ctx).Debug().Msgf("SSH Key %q already exists, using existing key", sshKey.Name)
			return k, nil
		}
		if sshKey.PublicKey != nil && *k.PublicKey == *sshKey.PublicKey {
			log.Ctx(ctx).Debug().Msgf("SSH Key %q already exists, using existing key", sshKey.Name)
			return k, nil
		}
	}
	// Key doesn't exist, create a new one

	log.Ctx(ctx).Debug().Msgf("Generating new SSH key %q", sshKey.Name)

	req := client.AddSSHKeyJSONRequestBody{
		Name:      sshKey.Name,
		PublicKey: sshKey.PublicKey,
	}
	res, err := s.client.AddSSHKeyWithResponse(ctx, req)
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

func (s *Session) findRegionForNode(ctx context.Context, nodeTypeID string) (string, error) {
	nodeTypes, err := s.ListNodeTypes(ctx, true)
	if err != nil {
		return "", fmt.Errorf("failed to list instance availability, err: %w", err)
	}

	for _, nt := range nodeTypes {
		if nt.ID == nodeTypeID {
			if len(nt.Regions) == 0 {
				continue
			}
			return nt.Regions[0], nil
		}
	}
	suggestion := ""
	b, err := json.Marshal(nodeTypes)
	if err != nil {
		log.Ctx(ctx).Warn().Msgf("Failed to marshal available instances to JSON: %v", err)
	} else {
		suggestion += string(b)
	}

	e := err503(fmt.Sprintf("No region with available capacity for node type %q", nodeTypeID), nil)
	e.Suggestion = suggestion
	return "", e
}

func (s *Session) HealthCheck(ctx context.Context) error {
	log.Ctx(ctx).Debug().Msg("Executing health check")
	_, err := s.ListNodeTypes(ctx, false)
	return err
}

func (s *Session) InitNode(ctx context.Context, sshKey types.SSHKey, nodeTypeID string, region *string) (types.Node, error) {
	log.Ctx(ctx).Debug().Msgf("Launching instance with SSH key %q", sshKey.Name)

	if region == nil {
		var err error
		var nr string
		nr, err = s.findRegionForNode(ctx, nodeTypeID)
		if err != nil {
			return types.Node{}, err
		}
		region = &nr
	}

	req := client.LaunchInstanceJSONRequestBody{
		FileSystemNames:  nil,
		InstanceTypeName: nodeTypeID,
		Name:             tools.Stringy("uw-" + random.GenerateRandomPhrase(3, "-")),
		Quantity:         tools.Inty(1),
		RegionName:       *region,
		SshKeyNames:      []string{sshKey.Name},
	}

	res, err := s.client.LaunchInstanceWithResponse(ctx, req)
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
			if strings.Contains(msg, "available capacity") {
				// Get a list of available instances
				instances, e := s.ListNodeTypes(ctx, true)
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
				err := err503(res.JSON400.Error.Message, nil)
				err.Suggestion = suggestion
				return types.Node{}, err
			}
			return types.Node{}, err400(res.JSON400.Error.Message, nil)
		}

		return types.Node{}, errUnknown(res.StatusCode(), err)
	}

	if len(res.JSON200.Data.InstanceIds) == 0 {
		return types.Node{}, fmt.Errorf("failed to launch instance")
	}

	return types.Node{
		ID:       res.JSON200.Data.InstanceIds[0],
		TypeID:   nodeTypeID,
		Region:   *region,
		KeyPair:  sshKey,
		Status:   types.StatusInitializing,
		Provider: types.LambdaLabsProvider,
	}, nil
}

func (s *Session) ListNodeTypes(ctx context.Context, filterAvailable bool) ([]types.NodeType, error) {
	log.Ctx(ctx).Debug().Msgf("Listing instance availability")

	res, err := s.client.InstanceTypesWithResponse(ctx)
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

		it := types.NodeType{
			ID:          id,
			Name:        &data.InstanceType.Description,
			Regions:     []string{},
			Price:       &data.InstanceType.PriceCentsPerHour,
			Description: nil,
			Provider:    types.LambdaLabsProvider,
			Specs: types.NodeSpecs{
				VCPUs:     data.InstanceType.Specs.Vcpus,
				Memory:    data.InstanceType.Specs.MemoryGib,
				GPUMemory: nil,
			},
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

func (s *Session) ListSSHKeys(ctx context.Context) ([]types.SSHKey, error) {
	log.Ctx(ctx).Debug().Msg("Listing SSH keys")

	res, err := s.client.ListSSHKeysWithResponse(ctx)
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

func (s *Session) NodeStatus(ctx context.Context, nodeID string) (types.SessionStatus, error) {
	log.Ctx(ctx).Debug().Msgf("Getting node status for node %q", nodeID)

	res, err := s.client.GetInstanceWithResponse(ctx, nodeID)
	if err != nil {
		return "", err
	}

	if res.JSON200 == nil {
		if res.JSON401 != nil {
			return "", err401(res.JSON401.Error.Message, nil)
		}
		if res.JSON403 != nil {
			return "", err403(res.JSON403.Error.Message, nil)
		}
		if res.JSON404 != nil {
			return "", err500(res.JSON404.Error.Message, fmt.Errorf("instance not found, %s", res.JSON404.Error.Message))
		}
		return "", errUnknown(res.StatusCode(), nil)
	}

	switch res.JSON200.Data.Status {
	case client.Active:
		return types.StatusRunning, nil
	case client.Booting:
		return types.StatusInitializing, nil
	case client.Terminated:
		return types.StatusTerminated, nil
	case client.Unhealthy:
		return types.StatusError, nil
	default:
		log.Ctx(ctx).Error().Msgf("Unknown lambda status %s", res.JSON200.Data.Status)
		return "", fmt.Errorf("unknown lambda status %s", res.JSON200.Data.Status)
	}
}

func (s *Session) TerminateNode(ctx context.Context, nodeID string) error {
	log.Ctx(ctx).Debug().Msgf("Terminating node %q", nodeID)

	req := client.TerminateInstanceJSONRequestBody{
		InstanceIds: []string{nodeID},
	}
	res, err := s.client.TerminateInstanceWithResponse(ctx, req)
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
		return errUnknown(res.StatusCode(), nil)
	}

	return nil
}

func (s *Session) Watch(ctx context.Context, nodeID string) (<-chan types.SessionStatus, <-chan error) {
	log.Ctx(ctx).Debug().Msgf("Watching node %q", nodeID)

	currentStatus := types.StatusInitializing
	statusch, errch := make(chan types.SessionStatus), make(chan error)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				status, e := s.NodeStatus(ctx, nodeID)
				if e != nil {
					errch <- fmt.Errorf("failed to get node state: %w", e)
				}
				if status == currentStatus {
					continue
				}
				currentStatus = status
				statusch <- status
			}
		}
	}()

	return statusch, errch
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

	return &Session{client: llClient}, nil
}
