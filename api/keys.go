package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/db"
)

// Context Keys should only be used inside the API package while parsing incoming requests
// either in the middleware or in the handlers. They should not be passed further into
// the call stack.
const (
	UserCtxKey    = "user"
	ProjectCtxKey = "project"
)

func getUserIDFromContext(ctx context.Context) uuid.UUID {
	uid, ok := ctx.Value(UserCtxKey).(uuid.UUID)
	if !ok {
		// This should never happen at runtime.
		log.Fatal().Msg("user not found in context")
	}
	return uid
}

func getProjectFromContext(ctx context.Context) *db.UnweaveProject {
	project, ok := ctx.Value(ProjectCtxKey).(db.UnweaveProject)
	if !ok {
		// This should never happen at runtime.
		log.Fatal().Msg("project not found in context")
	}
	return &project
}
