// Package runtime manages the lifecycle of a session.
package runtime

import (
	"context"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/builder"
	"github.com/unweave/unweave/vault"
	"github.com/unweave/unweave/wip/conductor/volume"
)

type Runtime struct {
	Node   Node
	Exec   Exec
	Volume volume.Provider
}

// Node represents an interactive session on a node. You can connect to it via SSH and
// train your ML models for example.
type Node interface {
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
	GetProvider() types.Provider
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
	// serve as a host to run containers that are accessible via SSH.
	InitNode(ctx context.Context, sshKey []types.SSHKey, spec types.HardwareSpec, region *string) (node types.Node, err error)
	// ListSSHKeys returns a list of all SSH keys associated with the provider.
	ListSSHKeys(ctx context.Context) ([]types.SSHKey, error)
	// ListNodeTypes returns a list of all node types available on the provider.
	ListNodeTypes(ctx context.Context, filterAvailable bool) ([]types.NodeType, error)
	// NodeStatus returns the status of the node running a session.
	NodeStatus(ctx context.Context, nodeID string) (types.Status, error)
	TerminateNode(ctx context.Context, nodeID string) error
	// Watch watches the status of the node.
	Watch(ctx context.Context, nodeID string) (<-chan types.Status, <-chan error)
}

type Exec interface {
	// Init initializes a new exec on a node. It creates environment the users code
	// will in with the provided build and configures ssh keys for interactive access.
	// If the persistentFS flag is set, the exec will be initialized with a persistent
	// file system. The call to Terminate is still required to handle the lifecycle of
	// the persistent file system. The flag just ensures the file system is capable of
	// being persisted.
	Init(ctx context.Context, node types.Node, config types.ExecConfig) (execID string, err error)
	// GetConnectionInfo returns the connection information for exec.
	GetConnectionInfo(ctx context.Context, execID string) (types.ConnectionInfo, error)
	// Terminate terminates a session.
	Terminate(ctx context.Context, execID string) error
	// Watch watches the status of the exec.
	Watch(ctx context.Context, execID string) (<-chan types.Status, <-chan error)
}

type Initializer interface {
	InitializeRuntime(ctx context.Context, userID string, provider types.Provider) (*Runtime, error)
	InitializeBuilder(ctx context.Context, userID string, builder string) (builder.Builder, error)
	InitializeVault(ctx context.Context) (vault.Vault, error)
}
