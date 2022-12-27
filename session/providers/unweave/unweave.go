package unweave

import (
	"github.com/unweave/unweave-v2/session/model"
)

type Runtime struct{}

func (r *Runtime) InitNode(model.SSHKey) (model.Node, error) {
	return model.Node{}, nil
}

func (r *Runtime) TerminateNode(nodeID string) error {
	return nil
}

func NewProvider(apiKey string) *Runtime {
	return &Runtime{}
}
