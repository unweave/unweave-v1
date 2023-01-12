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
	SessionCtxKey = "session"
)

func SetUserIDInContext(ctx context.Context, uid uuid.UUID) context.Context {
	return context.WithValue(ctx, UserCtxKey, uid)
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	uid, ok := ctx.Value(UserCtxKey).(uuid.UUID)
	if !ok {
		// This should never happen at runtime.
		log.Fatal().Msg("user not found in context")
	}
	return uid
}

func SetProjectInContext(ctx context.Context, project db.UnweaveProject) context.Context {
	return context.WithValue(ctx, ProjectCtxKey, project)
}

func GetProjectFromContext(ctx context.Context) *db.UnweaveProject {
	project, ok := ctx.Value(ProjectCtxKey).(db.UnweaveProject)
	if !ok {
		// This should never happen at runtime.
		log.Fatal().Msg("project not found in context")
	}
	return &project
}

func SetSessionInContext(ctx context.Context, session db.UnweaveSession) context.Context {
	return context.WithValue(ctx, SessionCtxKey, session)
}

func GetSessionFromContext(ctx context.Context) *db.UnweaveSession {
	session, ok := ctx.Value(SessionCtxKey).(db.UnweaveSession)
	if !ok {
		// This should never happen at runtime.
		log.Fatal().Msg("session not found in context")
	}
	return &session
}
