package lambdalabs

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type SessionRuntime struct{}

func (s *SessionRuntime) Init(ctx context.Context, node types.Node, sshKeys []types.SSHKey, image string) error {
	// noop - not implemented
	return nil
}

func (s *SessionRuntime) Exec(ctx context.Context, session string, execID string, params types.ExecCtx, isInteractive bool) error {
	// noop - not implemented
	return nil
}

func NewSessionRuntime(apiKey string) (*SessionRuntime, error) {
	return &SessionRuntime{}, nil
}
