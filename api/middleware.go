package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/db"
)

// withUserContext is a helper middleware that fakes an authenticated user. It should only
// be user for development or when self-hosting.
func withUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx,
			ContextKeyUser,
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withProjectContext is a helper middleware that parsed the project id from the url and
// verifies it exists in the db. It should only be user for development or when self-hosting.
func withProjectContext(dbq db.Querier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
			if err != nil {
				render.Render(w, r, &HTTPError{
					Code:       400,
					Message:    "Invalid project id",
					Suggestion: "Make sure the project id is a valid UUID",
				})
				return
			}

			project, err := dbq.ProjectGet(ctx, projectID)
			if err != nil {
				if err == sql.ErrNoRows {
					render.Render(w, r, &HTTPError{
						Code:       404,
						Message:    "Project not found",
						Suggestion: "Make sure the project id is valid",
					})
					return
				}
				log.Error().
					Err(err).
					Msgf("Error fetching session %q", projectID)

				render.Render(w, r, ErrInternalServer("Failed to terminate session"))
				return
			}

			ctx = context.WithValue(ctx, ContextKeyProject, project)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
