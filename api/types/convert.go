package types

import (
	"github.com/unweave/unweave/db"
)

func DBSessionStatusToAPIStatus(status db.UnweaveSessionStatus) SessionStatus {
	return SessionStatus(status)
}
