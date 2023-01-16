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
	"github.com/unweave/unweave/types"
)

// Context Keys should only be used inside the API package while parsing incoming requests
// either in the middleware or in the handlers. They should not be passed further into
// the call stack.
const (
	UserCtxKey    = "user"
	ProjectCtxKey = "project"
	SessionCtxKey = "session"
)

func SetSessionInContext(ctx context.Context, session types.Session) context.Context {
	return context.WithValue(ctx, SessionCtxKey, session)
}

func GetSessionFromContext(ctx context.Context) *types.Session {
	session, ok := ctx.Value(SessionCtxKey).(types.Session)
	if !ok {
		// This should never happen at runtime.
		log.Fatal().Msg("session not found in context")
	}
	return &session
}

// withSessionCtx is a helper middleware that parsed the session id from the url and
// verifies it exists in the db.
func withSessionCtx(store *Store) func(http.Handler) http.Handler {
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

			session, err := store.Session.Get(ctx, sessionID)
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
