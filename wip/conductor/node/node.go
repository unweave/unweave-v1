package node

import (
	"context"
	"io"

	"github.com/unweave/unweave-v1/wip/conductor/types"
)

type Node struct {
	ID      string
	State   types.NodeState
	Network types.Network
	Spec    types.Spec
}

// Driver is an interface that allows performing operations on a node. The implementing
// type should handle the details of networking protocols, different cloud providers etc.
type Driver interface {
	NodeRunCommand(ctx context.Context, id string, command []string, env []string) error
	NodeRunScript(ctx context.Context, id string, script, workingDir string, env []string) error
	NodeTransferFile(ctx context.Context, id string, src io.ReadCloser, dst string) error
}

// Provider is an interface that allows operating on VMs from different providers.
// AWS, GCP, DigitalOcean etc. are all examples of providers.
type Provider interface {
	NodeCreate(ctx context.Context) (string, error)
	NodeDelete(ctx context.Context, id string) error
	NodeInit(ctx context.Context, id string, options ...func(Driver)) (*Node, error)
	NodeList(ctx context.Context) ([]Node, error)
}
