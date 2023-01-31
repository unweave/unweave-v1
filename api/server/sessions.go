package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/runtime"
)

func SessionsCreate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		zerolog.Ctx(ctx).Info().Msgf("Executing SessionsCreate request")

		scr := types.SessionCreateParams{}
		if err := render.Bind(r, &scr); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Invalid request body"))
			return
		}

		userID := GetUserIDFromContext(ctx)
		project := GetProjectFromContext(ctx)
		srv := NewCtxService(rti, userID)

		session, err := srv.Session.Create(ctx, project.ID, scr)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to create session"))
			return
		}

		// TODO: watch status
		render.JSON(w, r, session)
	}
}

func SessionsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SessionsGet request")

		userID := GetUserIDFromContext(ctx)
		sess := GetSessionFromContext(ctx)
		srv := NewCtxService(rti, userID)

		session, err := srv.Session.Get(ctx, sess.ID)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to get session"))
			return
		}

		render.JSON(w, r, types.SessionGetResponse{Session: *session})
	}
}

func SessionsList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)
		project := GetProjectFromContext(ctx)

		log.Ctx(ctx).Info().Msgf("Executing SessionsList request")

		srv := NewCtxService(rti, userID)
		sessions, err := srv.Session.List(ctx, project.ID)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to list sessions"))
			return
		}
		render.JSON(w, r, types.SessionsListResponse{Sessions: sessions})
	}
}

func SessionsTerminate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)

		log.Ctx(ctx).
			Info().
			Msgf("Executing SessionsTerminate request for user %q", userID)

		session := GetSessionFromContext(ctx)
		srv := NewCtxService(rti, userID)

		params := types.SessionTerminateRequestParams{}
		render.Bind(r, &params)

		if err := srv.Session.Terminate(ctx, session.ID, params.ProviderToken); err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to terminate session"))
			return
		}
		render.Status(r, http.StatusOK)
	}
}
