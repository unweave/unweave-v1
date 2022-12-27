package lambdalabs

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/session/model"
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

type SSHKey struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

func (r *Runtime) InitNode(sshKey model.SSHKey) (model.Node, error) {
	// If the SSH key is not provided, generate a new one
	if sshKey.Name == "" {
		req := addSSHKeyRequest{
			SSHKey: SSHKey{
				Name:      "uw-generated-key",
				PublicKey: "",
			},
		}

		log.Info().
			Str(model.RuntimeProviderKey, model.LambdaLabsProvider.String()).
			Msgf("SSH Key not provided, generating new key")

		res, err := addSSHKey(req)
		if err != nil {
			return model.Node{}, err
		}

		sshKey.Name = res.Data.Name
		sshKey.PublicKey = res.Data.PublicKey
	}

	// Launch instance
	launchReq := launchInstanceRequest{
		RegionName:      string(r.Region),
		InstanceType:    string(r.Type),
		SSHKeyNames:     []string{sshKey.Name},
		FileSystemNames: nil,
		Quantity:        1,
		Name:            "uw-initialized-instance",
	}
	res, err := launchInstance(launchReq)
	if err != nil {
		return model.Node{}, fmt.Errorf("failed to launch instance, err: %v", err)
	}

	return model.Node{
		ID:      res.Data.InstanceIDs[0],
		KeyPair: sshKey,
		Status:  model.StatusInitializingNode,
	}, nil
}

func (r *Runtime) TerminateNode(nodeID string) error {
	return nil
}

func NewProvider(apiKey string) *Runtime {
	// Load LambdaLabsProvider credentials
	return &Runtime{}
}
