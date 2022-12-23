package unweave

import "github.com/unweave/unweave-v2/session/runtime"

type Runtime struct {
}

func (u *Runtime) InitNode() (runtime.SSHConnection, error) {
	return runtime.SSHConnection{}, nil
}

func (u *Runtime) StopNode() error {
	return nil
}

func NewProvider() Runtime {
	// Load credentials to make sure we're running on the Unweave platform
	return Runtime{}
}