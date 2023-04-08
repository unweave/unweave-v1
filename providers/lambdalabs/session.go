package lambdalabs

import (
	"context"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/tools/random"
)

type SessionRuntime struct{}

func (s *SessionRuntime) Init(ctx context.Context, node types.Node, sshKeys []types.SSHKey, image string) (string, error) {
	// noop - not implemented
	str, err := random.GenerateRandomString(11)
	if err != nil {
		return "", err
	}
	sessionID := "se_" + str
	return sessionID, nil
}

func (s *SessionRuntime) Exec(ctx context.Context, session string, execID string, params types.ExecCtx, isInteractive bool) error {
	// noop - not implemented
	return nil
}

func (s *SessionRuntime) Terminate(ctx context.Context, sessionID string) error {
	// noop - not implemented
	return nil
}

func NewSessionRuntime(apiKey string) (*SessionRuntime, error) {
	return &SessionRuntime{}, nil
}
