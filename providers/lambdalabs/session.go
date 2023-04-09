package lambdalabs

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
)

type ExecRuntime struct {
	node *NodeRuntime
}

func (e *ExecRuntime) Init(ctx context.Context, node types.Node, sshKeys []types.SSHKey, image string) (string, error) {
	// Session and Node ID are the same for LambdaLabs (for now)
	return node.ID, nil
}

func (e *ExecRuntime) Exec(ctx context.Context, session string, execID string, params types.ExecCtx, isInteractive bool) error {
	// noop - not implemented
	return nil
}

func (e *ExecRuntime) GetConnectionInfo(ctx context.Context, execID string) (types.ConnectionInfo, error) {
	return e.node.GetConnectionInfo(ctx, execID)
}

func (e *ExecRuntime) Terminate(ctx context.Context, sessionID string) error {
	// Session and Node ID are the same for LambdaLabs (for now)
	log.Ctx(ctx).Debug().Str("sessionID", sessionID).Msg("terminating session")
	return e.node.TerminateNode(ctx, sessionID)
}

func (e *ExecRuntime) Watch(ctx context.Context, execID string) (<-chan types.NodeStatus, <-chan error) {
	// Session and Node ID are the same for LambdaLabs (for now)
	log.Ctx(ctx).Debug().Str("execID", execID).Msg("watching exec")
	return e.node.Watch(ctx, execID)
}

func NewSessionRuntime(apiKey string) (*ExecRuntime, error) {
	nodeRuntime, err := NewNodeRuntime(apiKey)
	if err != nil {
		return nil, err
	}
	return &ExecRuntime{node: nodeRuntime}, nil

}
