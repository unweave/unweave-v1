package session

import "github.com/unweave/unweave-v2/session/model"

type Runtime interface {
	InitNode(sshKey model.SSHKey) (node model.Node, err error)
	TerminateNode(nodeID string) error
}
