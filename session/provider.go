package session

import "github.com/unweave/unweave-v2/types"

type Runtime interface {
	InitNode(sshKey types.SSHKey) (node types.Node, err error)
	TerminateNode(nodeID string) error
}
