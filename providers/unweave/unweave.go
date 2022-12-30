package unweave

import (
	"context"

	"github.com/unweave/unweave-v2/types"
)

type Session struct{}

func (r *Session) ListSSHKeys(ctx context.Context) ([]types.SSHKey, error) {
	return []types.SSHKey{}, nil
}

func (r *Session) AddSSHKey(context.Context, types.SSHKey) (types.SSHKey, error) {
	return types.SSHKey{}, nil
}

func (r *Session) InitNode(context.Context, types.SSHKey) (types.Node, error) {
	return types.Node{}, nil
}

func (r *Session) TerminateNode(ctx context.Context, nodeID string) error {
	return nil
}

func NewSessionProvider(apiKey string) (*Session, error) {
	return &Session{}, nil
}
