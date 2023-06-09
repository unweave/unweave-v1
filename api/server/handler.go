package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/middleware"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/wip/conductor/volume"
)

// Builder

// BuildsCreate expects a request body containing both the build context and the json
// params for the build.
//
//		eg. curl -X POST \
//				 -H 'Authorization: Bearer <token>' \
//	 		 	 -H 'Content-Type: multipart/form-data' \
//	 		 	 -F context=@context.zip \
//	 		 	 -F 'params={"builder": "docker"}'
//				 https://<api-host>/builds
func BuildsCreate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing BuildsCreate request")

		ibp := &types.BuildsCreateParams{}
		if err := ibp.Bind(r); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request body"))
			return
		}

		userID := middleware.GetUserIDFromContext(ctx)
		accountID := middleware.GetAccountIDFromContext(ctx)
		projectID := middleware.GetProjectIDFromContext(ctx)

		srv := NewCtxService(rti, accountID, userID)

		buildID, err := srv.Builder.Build(ctx, projectID, ibp)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to build image"))
			return
		}

		res := &types.BuildsCreateResponse{BuildID: buildID}
		render.JSON(w, r, res)
	}
}

// BuildsGet returns the details of a build. If the query param `logs` is set to
// true, the logs of the build will be returned as well.
func BuildsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing BuildsGet request")

		buildID := chi.URLParam(r, "buildID")
		getLogs := r.URL.Query().Get("logs") == "true"
		usedBy := r.URL.Query().Get("usedBy") == "true"

		userID := middleware.GetUserIDFromContext(ctx)
		accountID := middleware.GetAccountIDFromContext(ctx)

		srv := NewCtxService(rti, accountID, userID)

		// get build from db
		build, err := db.Q.BuildGet(ctx, buildID)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to get build"))
			return
		}

		var st *time.Time
		if build.StartedAt.Valid {
			st = &build.StartedAt.Time
		}
		var ft *time.Time
		if build.FinishedAt.Valid {
			ft = &build.FinishedAt.Time
		}

		res := &types.BuildsGetResponse{
			Build: types.Build{
				BuildID:     buildID,
				Name:        build.Name,
				ProjectID:   build.ProjectID,
				Status:      string(build.Status),
				BuilderType: build.BuilderType,
				CreatedAt:   build.CreatedAt,
				StartedAt:   st,
				FinishedAt:  ft,
			},
			UsedBySessions: []types.Exec{},
			Logs:           nil,
		}

		if usedBy {
			sessions, err := db.Q.BuildGetUsedBy(ctx, buildID)
			if err != nil {
				render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to get sessions using build"))
				return
			}

			ubs := make([]types.Exec, len(sessions))
			for i, s := range sessions {
				s := s
				ubs[i] = types.Exec{
					ID:   s.ID,
					Name: s.Name,
					//SSHKey:     types.SSHKey{}, // TODO: should we return this here?
					//Connection: nil,            // TODO: should we return this here?
					Status: types.Status(s.Status),
					//CreatedAt:  &s.CreatedAt,
					//NodeTypeID: s.NodeID,
					Region:   s.Region,
					Provider: types.Provider(s.Provider),
				}
			}
			res.UsedBySessions = ubs
		}

		if getLogs {
			logs, err := srv.Builder.GetLogs(ctx, buildID)
			if err != nil {
				render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to get build logs"))
				return
			}
			res.Logs = &logs
		}
		render.JSON(w, r, res)
	}
}

// Provider

// NodeTypesList returns a list of node types available for the user. If the query param
// `available` is set to true, only node types that are currently available to be
// scheduled will be returned.
func NodeTypesList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		provider := types.Provider(chi.URLParam(r, "provider"))
		log.Ctx(ctx).Info().Msgf("Executing NodeTypesList request for provider %s", provider)

		filterAvailable := r.URL.Query().Get("available") == "true"

		userID := middleware.GetUserIDFromContext(ctx)
		srv := NewCtxService(rti, "", userID)

		nodeTypes, err := srv.Provider.ListNodeTypes(ctx, provider, filterAvailable)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list node types"))
			return
		}

		res := &types.NodeTypesListResponse{NodeTypes: nodeTypes}
		render.JSON(w, r, res)
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
			render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request body"))
			return
		}

		userID := middleware.GetUserIDFromContext(ctx)
		accountID := middleware.GetAccountIDFromContext(ctx)

		srv := NewCtxService(rti, accountID, userID)

		name, err := srv.SSHKey.Add(ctx, params)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to add SSH key"))
			return
		}
		res := types.SSHKeyResponse{
			Name:       name,
			PublicKey:  params.PublicKey,
			PrivateKey: "",
		}
		render.JSON(w, r, &res)
	}
}

func SSHKeyGenerate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SSHKeyCreate request")

		params := types.SSHKeyGenerateParams{}
		if err := render.Bind(r, &params); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request body"))
			return
		}

		userID := middleware.GetUserIDFromContext(ctx)
		accountID := middleware.GetAccountIDFromContext(ctx)

		srv := NewCtxService(rti, accountID, userID)

		name, prv, pub, err := srv.SSHKey.Generate(ctx, params)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to generate SSH key"))
			return
		}

		res := types.SSHKeyResponse{
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

		userID := middleware.GetUserIDFromContext(ctx)
		accountID := middleware.GetAccountIDFromContext(ctx)

		srv := NewCtxService(rti, accountID, userID)

		keys, err := srv.SSHKey.List(ctx)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to list SSH keys"))
			return
		}

		res := types.SSHKeyListResponse{Keys: keys}
		render.JSON(w, r, res)
	}
}

func VolumeCreate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing VolumeCreate request")

		params := types.VolumeCreateParams{}
		if err := render.Bind(r, &params); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), types.ErrHTTPBadRequest(err, "Invalid request body"))
			return
		}

		accountID := middleware.GetAccountIDFromContext(ctx)
		projectID := middleware.GetProjectIDFromContext(ctx)
		volStore := GetVolumeStore()

		rt, err := rti.InitializeRuntime(ctx, accountID, params.Provider)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to initialize provider"))
			return
		}

		handler := volume.NewVolumeService(projectID, rt.Volume, volStore)

		vol, err := handler.Create(ctx, params.Size)
		if err != nil {
			render.Render(w, r.WithContext(ctx), types.ErrHTTPError(err, "Failed to create volume"))
			return
		}

		render.JSON(w, r, vol)
	}
}
