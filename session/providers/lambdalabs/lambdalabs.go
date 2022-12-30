package lambdalabs

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/pkg/random"
	"github.com/unweave/unweave-v2/types"
)

type InstanceDetails struct {
	Type   InstanceType `json:"type"`
	Region Region       `json:"region"`
	// TODO:
	// 	- Filesystems
}

type Runtime struct {
	InstanceDetails
}

func (r *Runtime) InitNode(sshKey types.SSHKey) (types.Node, error) {
	// If the SSH key is not provided, generate a new one
	if sshKey.Name == nil && sshKey.PublicKey == nil {
		name := "uw-generated-key-" + random.GenerateRandomPhrase(4, "-")
		req := addSSHKeyRequest{Name: name, PublicKey: ""}

		log.Info().
			Str(types.RuntimeProviderKey, types.LambdaLabsProvider.String()).
			Msgf("SSH Key not provided, generating new key")

		res, err := addSSHKey(req)
		if err != nil {
			return types.Node{}, err
		}

		sshKey.Name = &res.Data.Name
		sshKey.PublicKey = &res.Data.PublicKey
	}

	// Launch instance
	launchReq := launchInstanceRequest{
		RegionName:      string(r.Region),
		InstanceType:    string(r.Type),
		SSHKeyNames:     []string{*sshKey.Name},
		FileSystemNames: nil,
		Quantity:        1,
		Name:            "uw-initialized-instance",
	}
	res, err := launchInstance(launchReq)
	if err != nil {
		return types.Node{}, fmt.Errorf("failed to launch instance, err: %v", err)
	}

	return types.Node{
		ID:      res.Data.InstanceIDs[0],
		KeyPair: sshKey,
		Status:  types.StatusInitializingNode,
	}, nil
}

func (r *Runtime) TerminateNode(nodeID string) error {
	return nil
}

func NewProvider(apiKey string) *Runtime {
	// Load LambdaLabsProvider credentials
	return &Runtime{}
}
