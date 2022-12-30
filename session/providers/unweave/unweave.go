package unweave

import (
	"context"

	"github.com/unweave/unweave-v2/types"
)

type Runtime struct{}

func (r *Runtime) ListSSHKeys(ctx context.Context) ([]types.SSHKey, error) {
	return []types.SSHKey{}, nil
}

func (r *Runtime) AddSSHKey(context.Context, types.SSHKey) (types.SSHKey, error) {
	return types.SSHKey{}, nil
}

func (r *Runtime) InitNode(context.Context, types.SSHKey) (types.Node, error) {
	return types.Node{}, nil
}

func (r *Runtime) TerminateNode(ctx context.Context, nodeID string) error {
	return nil
}

func NewProvider(apiKey string) (*Runtime, error) {
	return &Runtime{}, nil
}
