// Package runtime manages the lifecycle of a session.
package runtime

import (
	"context"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/builder"
)

type Runtime struct {
	Session
}

// Session represents an interactive session on a node. You can connect to it via SSH and
// train your ML models for example.
type Session interface {
	// AddSSHKey adds a new SSH key to the provider.
	//
	// If the sshKey.PublicKey is nil, the provider will generate a new keypair with the
	// given sshKey.Name and return both the public and private keys in the response.
	// Otherwise, if the shhKey.PublicKey is not nil, the provider will verify that the
	// key is valid and return the public key and name in the response.
	//
	// sshKey.Name must not be empty.
	// This operation must be idempotent. i.e. if the sshKey.Name or sshKey.PublicKey
	// already exist with the provider, this should be a no-op. In this case, both the
	// name and public key should match those with the provider.
	AddSSHKey(ctx context.Context, sshKey types.SSHKey) (types.SSHKey, error)
	// GetProvider returns the provider.
	GetProvider() types.RuntimeProvider
	// GetConnectionInfo returns the connection information for the node running a session.
	GetConnectionInfo(ctx context.Context, nodeID string) (types.ConnectionInfo, error)
	// HealthCheck performs a health check on the provider.
	HealthCheck(ctx context.Context) error
	// InitNode initializes a new node on the provider.
	//
	// It should automatically select the most appropriate region if one is not specified.
	// The implementation should choose the level of abstraction this method is
	// implemented at. For example, it could be implemented at a VM level for a bare-metal
	// provider, at a container level, batch job level, etc. In each case, the node must
	// be accessible via SSH.
	InitNode(ctx context.Context, sshKey types.SSHKey, nodeTypeID string, region *string) (node types.Node, err error)
	// ListSSHKeys returns a list of all SSH keys associated with the provider.
	ListSSHKeys(ctx context.Context) ([]types.SSHKey, error)
	// ListNodeTypes returns a list of all node types available on the provider.
	ListNodeTypes(ctx context.Context, filterAvailable bool) ([]types.NodeType, error)
	// NodeStatus returns the status of the node running a session.
	NodeStatus(ctx context.Context, nodeID string) (types.SessionStatus, error)
	TerminateNode(ctx context.Context, nodeID string) error
	// Watch watches the status of the node running a session.
	Watch(ctx context.Context, nodeID string) (<-chan types.SessionStatus, <-chan error)
}

type Initializer interface {
	InitializeRuntime(ctx context.Context, accountID string, provider types.RuntimeProvider) (*Runtime, error)
	InitializeBuilder(ctx context.Context, accountID string, builder string) (builder.Builder, error)
}
