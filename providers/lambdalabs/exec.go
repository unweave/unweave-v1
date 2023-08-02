package lambdalabs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/providers/lambdalabs/client"
	"github.com/unweave/unweave-v1/services/execsrv"
	"github.com/unweave/unweave-v1/tools"
	"github.com/unweave/unweave-v1/tools/random"
)

func (d *Driver) ExecCreate(ctx context.Context, project, image string, spec types.HardwareSpec, network types.ExecNetwork, volumes []types.ExecVolume, pubKeys []string, region *string) (string, error) {
	if len(pubKeys) == 0 {
		return "", fmt.Errorf("no ssh keys provided")
	}

	kayNames := make([]string, len(pubKeys))

	for idx, pubKey := range pubKeys {
		name, err := d.sshKeyRegister(ctx, pubKey)
		if err != nil {
			return "", fmt.Errorf("failed to register ssh key, err: %v", err)
		}
		kayNames[idx] = name
	}

	var nodeTypeID = spec.GPU.Type

	if region == nil {
		var err error
		var nr string
		nr, err = d.findRegionForNode(ctx, nodeTypeID)
		if err != nil {
			return "", err
		}
		region = &nr
	}

	req := client.LaunchInstanceJSONRequestBody{
		FileSystemNames:  nil,
		InstanceTypeName: nodeTypeID,
		Name:             tools.Stringy("uw-" + random.GenerateRandomPhrase(3, "-")),
		Quantity:         tools.Inty(1),
		RegionName:       *region,
		SshKeyNames:      kayNames,
	}

	res, err := d.client.LaunchInstanceWithResponse(ctx, req)
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
		if res.JSON500 != nil {
			return "", err500(res.JSON500.Error.Message, nil)
		}
		if res.JSON404 != nil {
			return "", err404(res.JSON404.Error.Message, nil)
		}

		// We get a 400 if the instance type is not available. We check for the available
		// instances and return them in the error message. Since this is not critical, we
		// can ignore if there are any errors in the process.
		if res.JSON400 != nil {
			suggestion := ""
			msg := strings.ToLower(res.JSON400.Error.Message)
			if strings.Contains(msg, "available capacity") {
				// Get a list of available instances
				instances, e := d.listNodeTypes(ctx, true)
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
				return "", err
			}
			return "", err400(res.JSON400.Error.Message, nil)
		}

		return "", errUnknown(res.StatusCode(), err)
	}

	if len(res.JSON200.Data.InstanceIds) == 0 {
		return "", fmt.Errorf("failed to launch instance")
	}

	return res.JSON200.Data.InstanceIds[0], nil
}

func (d *Driver) ExecDriverName() string {
	return "lambdalabs"
}

func (d *Driver) listNodeTypes(ctx context.Context, filterAvailable bool) ([]types.NodeType, error) {
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

func (d *Driver) findRegionForNode(ctx context.Context, nodeTypeID string) (string, error) {
	nodeTypes, err := d.listNodeTypes(ctx, true)
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

func (d *Driver) Get(ctx context.Context, id string) (types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) ExecGetStatus(ctx context.Context, execID string) (types.Status, error) {
	res, err := d.client.GetInstanceWithResponse(ctx, execID)
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

func (d *Driver) List(ctx context.Context, project string) ([]types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) ExecProvider() types.Provider {
	return types.LambdaLabsProvider
}

func (d *Driver) sshKeyList(ctx context.Context) ([]types.SSHKey, error) {
	res, err := d.client.ListSSHKeysWithResponse(ctx)
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

func (d *Driver) sshKeyRegister(ctx context.Context, pubKey string) (string, error) {
	keys, err := d.sshKeyList(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list ssh keys, err: %w", err)
	}

	for _, k := range keys {
		if *k.PublicKey == pubKey {
			log.Ctx(ctx).Debug().Msgf("SSH Key already exists, using existing key")
			return k.Name, nil
		}
	}

	// Key doesn't exist, create a new one

	name := "uw:" + random.GenerateRandomPhrase(4, "-")
	log.Ctx(ctx).Debug().Msgf("Generating new SSH key %q", name)

	req := client.AddSSHKeyJSONRequestBody{
		Name:      name,
		PublicKey: &pubKey,
	}
	res, err := d.client.AddSSHKeyWithResponse(ctx, req)
	if err != nil {
		return "", err
	}
	if res.JSON200 == nil {
		err = fmt.Errorf("failed to generate SSH key")
		if res.JSON401 != nil {
			return "", err401(res.JSON401.Error.Message, err)
		}
		if res.JSON403 != nil {
			return "", err403(res.JSON403.Error.Message, err)
		}
		if res.JSON400 != nil {
			return "", err400(res.JSON400.Error.Message, err)
		}
		return "", errUnknown(res.StatusCode(), err)
	}

	return name, nil
}

func (d *Driver) ExecSpec(ctx context.Context, execID string) (types.HardwareSpec, error) {
	instance, err := d.client.GetInstanceWithResponse(ctx, execID)
	if err != nil {
		return types.HardwareSpec{}, &types.Error{
			Code:     http.StatusInternalServerError,
			Message:  "Failed to make request to LambdaLabs API",
			Provider: types.LambdaLabsProvider,
			Err:      fmt.Errorf("failed to get instance, err: %w", err),
		}
	}

	if instance.JSON200 == nil {
		err = fmt.Errorf("failed to get instance")
		if instance.JSON401 != nil {
			return types.HardwareSpec{}, err401(instance.JSON401.Error.Message, err)
		}
		if instance.JSON403 != nil {
			return types.HardwareSpec{}, err403(instance.JSON403.Error.Message, err)
		}
		if instance.JSON404 != nil {
			return types.HardwareSpec{}, err404(instance.JSON404.Error.Message, err)
		}
		return types.HardwareSpec{}, errUnknown(instance.StatusCode(), err)
	}

	gpuCount, err := parseGPUCount(instance.JSON200.Data.InstanceType.Name)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse number of GPUs")
		gpuCount = 0
	}

	gpuMem, err := parseGPUMemory(instance.JSON200.Data.InstanceType.Name)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse GPU memory")
		gpuMem = 0
	}

	spec := getHardwareSpecFromInstance(instance.JSON200.Data, gpuCount, gpuMem)
	return spec, nil
}

func (d *Driver) ExecStats(_ context.Context, _ string) (execsrv.Stats, error) {
	panic("implement me")
}

func (d *Driver) ExecTerminate(_ context.Context, _ string) error {
	panic("implement me")
}

func (d *Driver) ExecConnectionInfo(_ context.Context, _ string) (types.ConnectionInfo, error) {
	panic("not implemented")
}
