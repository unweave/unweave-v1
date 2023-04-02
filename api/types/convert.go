package types

import (
	"github.com/unweave/unweave/db"
)

func DBSessionStatusToAPIStatus(status db.UnweaveSessionStatus) NodeStatus {
	return NodeStatus(status)
}
