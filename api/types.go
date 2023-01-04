package api

import (
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/types"
)

func dbSessionStatusToAPIStatus(status db.UnweaveSessionStatus) types.SessionStatus {
	return types.SessionStatus(status)
}
