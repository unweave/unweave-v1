package unweave

import "github.com/unweave/unweave-v2/types"

type Runtime struct{}

func (r *Runtime) InitNode(types.SSHKey) (types.Node, error) {
	return types.Node{}, nil
}

func (r *Runtime) TerminateNode(nodeID string) error {
	return nil
}

func NewProvider(apiKey string) *Runtime {
	return &Runtime{}
}
