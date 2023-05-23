package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

func GetProviderIDFromProject(projectID string) (string, error) {
	ctx := context.Background()

	_, err := db.Q.ProjectGet(ctx, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &types.Error{
				Code:    http.StatusNotFound,
				Message: "Project not found",
			}
		}
	}

	return "", nil
}
