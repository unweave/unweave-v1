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
	AccountIDCtxKey     = "accountID"
	BuildIDCtxKey       = "buildID"
	ProjectIDCtxKey     = "project"
	SessionIDCtxKey     = "session"
	SessionStatusCtxKey = "sessionStatus"
)

func SetAccountIDInContext(ctx context.Context, aid string) context.Context {
	return context.WithValue(ctx, AccountIDCtxKey, aid)
}

func GetAccountIDFromContext(ctx context.Context) string {
	uid, ok := ctx.Value(AccountIDCtxKey).(string)
	if !ok || uid == "" {
		// This should never happen at runtime.
		log.Error().Msg("account not found in context")
		panic("account not found in context")
	}
	return uid
}

func SetProjectIDInContext(ctx context.Context, projectID string) context.Context {
	return context.WithValue(ctx, ProjectIDCtxKey, projectID)
}

func GetProjectIDFromContext(ctx context.Context) string {
	projectID, ok := ctx.Value(ProjectIDCtxKey).(string)
	if !ok || projectID == "" {
		// This should never happen at runtime.
		log.Error().Msg("project not found in context")
		panic("project not found in context")
	}
	return projectID
}

func SetSessionIDInContext(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDCtxKey, sessionID)
}

func GetSessionIDFromContext(ctx context.Context) string {
	sessionID, ok := ctx.Value(SessionIDCtxKey).(string)
	if !ok || sessionID == "" {
		// This should never happen at runtime.
		log.Error().Msg("session not found in context")
		panic("session not found in context")
	}
	return sessionID
}

// withAccountCtx is a helper middleware that fakes an authenticated account. It should only
// be user for development or when self-hosting.
func withAccountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accountID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx = context.WithValue(ctx, AccountIDCtxKey, accountID)
		ctx = log.With().Stringer(AccountIDCtxKey, accountID).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withProjectCtx is a helper middleware that parsed the project id from the url and
// verifies it exists in the db.
func withProjectCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := chi.URLParam(r, "projectID")

		_, err := db.Q.ProjectGet(ctx, projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &types.Error{
					Code:       http.StatusNotFound,
					Message:    "Project not found",
					Suggestion: "Make sure the project id is valid",
				})
				return
			}

			err = fmt.Errorf("failed to fetch project from db %q: %w", projectID, err)
			render.Render(w, r.WithContext(ctx),
				ErrInternalServer(err, "Failed to fetch project"))
			return
		}

		ctx = context.WithValue(ctx, ProjectIDCtxKey, projectID)
		ctx = log.With().Str(ProjectIDCtxKey, projectID).Logger().WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withSessionCtx is a helper middleware that parsed the session id from the url and
// verifies it exists in the db.
func withSessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionID := chi.URLParam(r, "sessionID")

		session, err := db.Q.SessionGet(ctx, sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &types.Error{
					Code:       http.StatusNotFound,
					Message:    "Session not found",
					Suggestion: "Make sure the session id is valid",
				})
				return
			}

			err = fmt.Errorf("failed to fetch session from db %q: %w", sessionID, err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, "Failed to fetch session"))
			return
		}

		ctx = context.WithValue(ctx, SessionIDCtxKey, session)
		ctx = log.With().Str(SessionIDCtxKey, session.ID).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
