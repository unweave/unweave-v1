package types

import (
	"github.com/unweave/unweave/db"
)

func DBSessionStatusToAPIStatus(status db.UnweaveExecStatus) Status {
	return Status(status)
}
