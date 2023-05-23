package node

import (
	"context"

	"github.com/unweave/unweave/wip/conductor/types"
)

type Node struct {
	ID    string
	State types.NodeState
	Spec  types.Spec
}

type Provider interface {
	NodeCreate(ctx context.Context) (string, error)
	NodeDelete(ctx context.Context, id string) error
	NodeInit(ctx context.Context, id string) error
	NodeList(ctx context.Context) ([]string, error)
}
