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
	ProjectIDCtxKey     = "project"
	SessionIDCtxKey     = "session"
	SessionStatusCtxKey = "sessionStatus"
)

func SetAccountIDInContext(ctx context.Context, uid uuid.UUID) context.Context {
	return context.WithValue(ctx, AccountIDCtxKey, uid)
}

func GetAccountIDFromContext(ctx context.Context) uuid.UUID {
	uid, ok := ctx.Value(AccountIDCtxKey).(uuid.UUID)
	if !ok {
		// This should never happen at runtime.
		log.Error().Msg("account not found in context")
		panic("account not found in context")
	}
	return uid
}

func SetProjectIDInContext(ctx context.Context, projectID uuid.UUID) context.Context {
	return context.WithValue(ctx, ProjectIDCtxKey, projectID)
}

func GetProjectIDFromContext(ctx context.Context) uuid.UUID {
	projectID, ok := ctx.Value(ProjectIDCtxKey).(uuid.UUID)
	if !ok {
		// This should never happen at runtime.
		log.Error().Msg("project not found in context")
		panic("project not found in context")
	}
	return projectID
}

func SetSessionIDInContext(ctx context.Context, sessionID uuid.UUID) context.Context {
	return context.WithValue(ctx, SessionIDCtxKey, sessionID)
}

func GetSessionIDFromContext(ctx context.Context) uuid.UUID {
	sessionID, ok := ctx.Value(SessionIDCtxKey).(uuid.UUID)
	if !ok {
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
		projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
		if err != nil {
			render.Render(w, r.WithContext(ctx), &types.Error{
				Code:       http.StatusBadRequest,
				Message:    "Invalid project id",
				Suggestion: "Make sure the project id is a valid UUID",
			})
			return
		}

		project, err := db.Q.ProjectGet(ctx, projectID)
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
				ErrInternalServer(err, "Failed to terminate session"))
			return
		}

		ctx = context.WithValue(ctx, ProjectIDCtxKey, project)
		ctx = log.With().Stringer(ProjectIDCtxKey, project.ID).Logger().WithContext(ctx)

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
			render.Render(w, r.WithContext(ctx), &types.Error{
				Code:       http.StatusBadRequest,
				Message:    "Invalid session id",
				Suggestion: "Make sure the session id is a valid UUID",
			})
			return
		}

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
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, "Failed to terminate session"))
			return
		}

		ctx = context.WithValue(ctx, SessionIDCtxKey, session)
		ctx = log.With().Stringer(SessionIDCtxKey, session.ID).Logger().WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
