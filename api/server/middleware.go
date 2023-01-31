package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

// Context Keys should only be used inside the API package while parsing incoming requests
// either in the middleware or in the handlers. They should not be passed further into
// the call stack.
const (
	UserIDCtxKey  = "user"
	ProjectCtxKey = "project"
	SessionCtxKey = "session"
)

func SetUserIDInContext(ctx context.Context, uid uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDCtxKey, uid)
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	uid, ok := ctx.Value(UserIDCtxKey).(uuid.UUID)
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

// withUserCtx is a helper middleware that fakes an authenticated user. It should only
// be user for development or when self-hosting.
func withUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx = context.WithValue(ctx, UserIDCtxKey, userID)
		ctx = log.With().Stringer(UserIDCtxKey, userID).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withProjectCtx is a helper middleware that parsed the project id from the url and
// verifies it exists in the db.
func withProjectCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
		if err != nil {
			render.Render(w, r.WithContext(ctx), &types.HTTPError{
				Code:       http.StatusBadRequest,
				Message:    "Invalid project id",
				Suggestion: "Make sure the project id is a valid UUID",
			})
			return
		}

		project, err := db.Q.ProjectGet(ctx, projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &types.HTTPError{
					Code:       http.StatusNotFound,
					Message:    "Project not found",
					Suggestion: "Make sure the project id is valid",
				})
				return
			}

			err = fmt.Errorf("failed to fetch project from db %q: %w", projectID, err)
			render.Render(w, r.WithContext(ctx),
				ErrInternalServer(err, "Failed to terminate session"))
			return
		}

		ctx = context.WithValue(ctx, ProjectCtxKey, project)
		ctx = log.With().Stringer(ProjectCtxKey, project.ID).Logger().WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withSessionCtx is a helper middleware that parsed the session id from the url and
// verifies it exists in the db.
func withSessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
		if err != nil {
			render.Render(w, r.WithContext(ctx), &types.HTTPError{
				Code:       http.StatusBadRequest,
				Message:    "Invalid session id",
				Suggestion: "Make sure the session id is a valid UUID",
			})
			return
		}

		session, err := db.Q.SessionGet(ctx, sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &types.HTTPError{
					Code:       http.StatusNotFound,
					Message:    "Session not found",
					Suggestion: "Make sure the session id is valid",
				})
				return
			}

			err = fmt.Errorf("failed to fetch session from db %q: %w", sessionID, err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, "Failed to terminate session"))
			return
		}

		ctx = context.WithValue(ctx, SessionCtxKey, session)
		ctx = log.With().Stringer(SessionCtxKey, session.ID).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
