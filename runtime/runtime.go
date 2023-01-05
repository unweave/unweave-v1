// Package runtime manages the lifecycle of a session.
package runtime

import (
	"context"

	"github.com/unweave/unweave/types"
)

type Runtime struct {
	Session
}

// Session represents an interactive session on a node. You can connect to it via SSH and
// train your ML models for example.
type Session interface {
	// AddSSHKey adds a new SSH key to the provider.
	//
	// If sshKey.Name is nil, the provider automatically generates a random, unique,
	// human-readable name. If the sshKey.PublicKey is nil, the provider will generate a
	// new keypair and return both the public and private keys in the response.
	// Otherwise, if the shhKey.PublicKey is not nil, the provider will verify that the
	// key is valid and return the public key and name in the response.
	AddSSHKey(ctx context.Context, sshKey types.SSHKey) (types.SSHKey, error)
	GetProvider() types.RuntimeProvider
	InitNode(ctx context.Context, sshKey types.SSHKey) (node types.Node, err error)
	// ListSSHKeys returns a list of all SSH keys associated with the provider.
	ListSSHKeys(ctx context.Context) ([]types.SSHKey, error)
	TerminateNode(ctx context.Context, nodeID string) error
}
