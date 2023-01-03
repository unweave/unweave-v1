package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/types"
)

func SessionsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		res := &types.Session{ID: id}
		render.JSON(w, r, res)
	}
}

func SessionsList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := []*types.Session{
			{ID: "1"},
		}
		render.JSON(w, r, res)
	}
}

type SessionTerminateResponse struct {
	Success bool `json:"success"`
}

func SessionsTerminate(rti runtime.Initializer, dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := getUserIDFromContext(ctx)

		logger := log.With().
			Str(ContextKeyUser, userID.String()).
			Logger()

		logger.Info().
			Msgf("Executing SessionsTerminate request for user %q", userID)

		// fetch from url params and try converting to uuid
		id := chi.URLParam(r, "sessionID")
		sessionID, err := uuid.Parse(id)
		if err != nil {
			render.Render(w, r, &HTTPError{
				Code:       400,
				Message:    "Invalid session id",
				Suggestion: "Make sure the session id is a valid UUID",
			})
			return
		}

		sess, err := dbq.SessionGet(ctx, sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, &HTTPError{
					Code:       404,
					Message:    "Session not found",
					Suggestion: "Make sure the session id is valid",
				})
				return
			}
			logger.Error().
				Err(err).
				Msgf("Error fetching session %q", sessionID)

			render.Render(w, r, ErrInternalServer("Failed to terminate session"))
			return
		}

		rt, err := rti.FromUser(sessionID, types.RuntimeProvider(sess.Runtime))
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to create runtime" + sess.Runtime)

			render.Render(w, r, ErrInternalServer("Failed to initialize runtime"))
			return
		}

		if err = rt.TerminateNode(ctx, sess.NodeID); err != nil {
			render.Render(w, r, ErrHTTPError(err, "Failed to terminate node"))
			return
		}
		if err = dbq.SessionSetTerminated(ctx, sessionID); err != nil {
			logger.Error().
				Err(err).
				Msgf("Failed to set session %q as terminated", sessionID)
		}

		render.JSON(w, r, &SessionTerminateResponse{Success: true})
	}
}
