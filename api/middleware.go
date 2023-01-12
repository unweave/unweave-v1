package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/db"
)

// withUserCtx is a helper middleware that fakes an authenticated user. It should only
// be user for development or when self-hosting.
func withUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := log.With().Logger().WithContext(r.Context())
		ctx = context.WithValue(ctx,
			UserCtxKey,
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withProjectCtx is a helper middleware that parsed the project id from the url and
// verifies it exists in the db.
func withProjectCtx(dbq db.Querier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.With().Logger().WithContext(r.Context())
			projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
			if err != nil {
				render.Render(w, r.WithContext(ctx), &HTTPError{
					Code:       http.StatusBadRequest,
					Message:    "Invalid project id",
					Suggestion: "Make sure the project id is a valid UUID",
				})
				return
			}

			project, err := dbq.ProjectGet(ctx, projectID)
			if err != nil {
				if err == sql.ErrNoRows {
					render.Render(w, r.WithContext(ctx), &HTTPError{
						Code:       http.StatusNotFound,
						Message:    "Project not found",
						Suggestion: "Make sure the project id is valid",
					})
					return
				}

				err = fmt.Errorf("failed to fetch project from db %q: %w", projectID, err)
				render.Render(w, r.WithContext(ctx), ErrInternalServer(err, "Failed to terminate session"))
				return
			}

			ctx = context.WithValue(ctx, ProjectCtxKey, project)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// withSessionCtx is a helper middleware that parsed the session id from the url and
// verifies it exists in the db.
func withSessionCtx(dbq db.Querier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.With().Logger().WithContext(r.Context())
			sessionID, err := uuid.Parse(chi.URLParam(r, "sessionID"))
			if err != nil {
				render.Render(w, r.WithContext(ctx), &HTTPError{
					Code:       http.StatusBadRequest,
					Message:    "Invalid session id",
					Suggestion: "Make sure the session id is a valid UUID",
				})
				return
			}

			session, err := dbq.SessionGet(ctx, sessionID)
			if err != nil {
				if err == sql.ErrNoRows {
					render.Render(w, r.WithContext(ctx), &HTTPError{
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
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
