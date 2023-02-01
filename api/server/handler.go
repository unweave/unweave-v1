package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/runtime"
)

// Provider

// NodeTypesList returns a list of node types available for the user
func NodeTypesList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		provider := types.RuntimeProvider(chi.URLParam(r, "provider"))
		log.Ctx(ctx).Info().Msgf("Executing NodeTypesList request for provider %s", provider)

		userID := GetUserIDFromContext(ctx)

		srv := NewCtxService(rti, userID)
		nodeTypes, err := srv.Provider.ListNodeTypes(ctx, provider)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to list node types"))
			return
		}

		res := &types.NodeTypesListResponse{NodeTypes: nodeTypes}
		render.JSON(w, r, res)
	}
}

// Sessions

func SessionsCreate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SessionsCreate request")

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

// SSH Keys

// SSHKeyAdd adds an SSH key to the user's account.
//
// This does not add the key to the user's configured providers. That is done lazily
// when the user first tries to use the key.
func SSHKeyAdd(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SSHKeyAdd request")

		params := types.SSHKeyAddParams{}
		if err := render.Bind(r, &params); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Invalid request body"))
			return
		}

		userID := GetUserIDFromContext(ctx)
		srv := NewCtxService(rti, userID)

		if err := srv.SSHKey.Add(ctx, params); err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to add SSH key"))
			return
		}
		render.JSON(w, r, &types.SSHKeyAddResponse{Success: true})
	}
}

func SSHKeyGenerate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SSHKeyCreate request")

		params := types.SSHKeyGenerateParams{}
		render.Bind(r, &params)

		userID := GetUserIDFromContext(ctx)
		srv := NewCtxService(rti, userID)

		name, prv, pub, err := srv.SSHKey.Generate(ctx, params)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to generate SSH key"))
			return
		}

		res := types.SSHKeyGenerateResponse{
			Name:       name,
			PublicKey:  pub,
			PrivateKey: prv,
		}
		render.JSON(w, r, &res)
	}
}

func SSHKeyList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SSHKeyList request")

		userID := GetUserIDFromContext(ctx)
		srv := NewCtxService(rti, userID)

		keys, err := srv.SSHKey.List(ctx)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to list SSH keys"))
			return
		}

		res := types.SSHKeyListResponse{Keys: keys}
		render.JSON(w, r, res)
	}
}
