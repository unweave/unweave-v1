package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/runtime"
)

// Builder

func ImagesBuild(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing ImagesBuild request")

		ibp := types.ImageBuildParams{}
		if err := render.Bind(r, &ibp); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Invalid request body"))
			return
		}

		accountID := GetAccountIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

		buildID, err := srv.Builder.Build(ctx, nil)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to build image"))
			return
		}

		res := &types.ImageBuildResponse{BuildID: buildID}
		render.JSON(w, r, res)
	}
}

// Provider

// NodeTypesList returns a list of node types available for the user
func NodeTypesList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		provider := types.RuntimeProvider(chi.URLParam(r, "provider"))
		log.Ctx(ctx).Info().Msgf("Executing NodeTypesList request for provider %s", provider)

		filterAvailable := r.URL.Query().Get("available") == "true"

		accountID := GetAccountIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

		nodeTypes, err := srv.Provider.ListNodeTypes(ctx, provider, filterAvailable)
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

		accountID := GetAccountIDFromContext(ctx)
		projectID := GetProjectIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

		session, err := srv.Session.Create(ctx, projectID, scr)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to create session"))
			return
		}

		go func() {
			c := context.Background()
			c = log.With().
				Stringer(AccountIDCtxKey, accountID).
				Str(ProjectIDCtxKey, projectID).
				Str(SessionIDCtxKey, session.ID).
				Logger().WithContext(c)

			if e := srv.Session.Watch(c, session.ID); e != nil {
				log.Ctx(ctx).Error().Err(e).Msgf("Failed to watch session")
			}
		}()

		render.JSON(w, r, session)
	}
}

func SessionsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SessionsGet request")

		accountID := GetAccountIDFromContext(ctx)
		sessionID := GetSessionIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

		session, err := srv.Session.Get(ctx, sessionID)
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
		accountID := GetAccountIDFromContext(ctx)
		projectID := GetProjectIDFromContext(ctx)
		listTerminated := r.URL.Query().Get("terminated") == "true"

		log.Ctx(ctx).Info().Msgf("Executing SessionsList request")

		srv := NewCtxService(rti, accountID)
		sessions, err := srv.Session.List(ctx, projectID, listTerminated)
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
		accountID := GetAccountIDFromContext(ctx)

		log.Ctx(ctx).
			Info().
			Msgf("Executing SessionsTerminate request")

		sessionID := GetSessionIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

		if err := srv.Session.Terminate(ctx, sessionID); err != nil {
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

		accountID := GetAccountIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

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

		accountID := GetAccountIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

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

		accountID := GetAccountIDFromContext(ctx)
		srv := NewCtxService(rti, accountID)

		keys, err := srv.SSHKey.List(ctx)
		if err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to list SSH keys"))
			return
		}

		res := types.SSHKeyListResponse{Keys: keys}
		render.JSON(w, r, res)
	}
}
