package lambdalabs

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/model"
	"github.com/unweave/unweave-v2/session/runtime"
)

type InstanceDetails struct {
	Type   InstanceType `json:"type"`
	Region Region       `json:"region"`
	// TODO:
	// 	- Filesystems
}

type LlRuntime struct {
	SSHKey
	InstanceDetails
}

type SSHKey struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

func (l *LlRuntime) InitNode() (runtime.Node, error) {
	node := runtime.Node{
		ID: "",
		KeyPair: runtime.SSHKeyPair{
			PublicKey: l.PublicKey,
		},
		Status: runtime.StatusInitializingNode,
	}

	// If the SSH key is not provided, generate a new one
	if l.SSHKey.Name == "" {
		req := AddSSHKeyRequest{
			SSHKey: SSHKey{
				Name:      "uw-generated-key",
				PublicKey: "",
			},
		}

		log.Info().
			Str(model.RuntimeProviderKey, model.LambdaLabsProvider.String()).
			Msgf("SSH Key not provided, generating new key")

		res, err := AddSSHKey(req)
		if err != nil {
			return runtime.Node{}, err
		}

		l.SSHKey.Name = res.Data.Name
		l.SSHKey.PublicKey = res.Data.PublicKey
		node.KeyPair.PrivateKey = res.Data.PrivateKey
		node.KeyPair.PublicKey = res.Data.PublicKey
	}

	// Launch instance
	launchReq := LaunchInstanceRequest{
		RegionName:      string(l.Region),
		InstanceType:    string(l.Type),
		SSHKeyNames:     []string{l.SSHKey.Name},
		FileSystemNames: nil,
		Quantity:        1,
		Name:            "uw-initialized-instance",
	}
	res, err := LaunchInstance(launchReq)
	if err != nil {
		return runtime.Node{}, fmt.Errorf("failed to launch instance, err: %v", err)
	}
	node.ID = res.Data.InstanceIDs[0]

	return node, nil
}

func (l *LlRuntime) TerminateNode() error {
	return nil
}

func NewProvider(key SSHKey) LlRuntime {
	// Load LambdaLabsProvider credentials
	return LlRuntime{
		SSHKey: key,
	}
}
