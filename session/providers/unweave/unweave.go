package unweave

import "github.com/unweave/unweave-v2/session/runtime"

type Runtime struct{}

func (u *Runtime) InitNode() (runtime.Node, error) {
	return "", nil
}

func (u *Runtime) TerminateNode() error {
	return nil
}

func NewProvider() Runtime {
	// Load credentials to make sure we're running on the Unweave platform

	// init node
	//
	return Runtime{}
}
