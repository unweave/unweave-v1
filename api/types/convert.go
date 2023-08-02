package types

import (
	"github.com/unweave/unweave-v1/db"
)

func DBSessionStatusToAPIStatus(status db.UnweaveExecStatus) Status {
	return Status(status)
}
