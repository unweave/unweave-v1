package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/runtime"
	"github.com/unweave/unweave-v2/types"
)

type SessionCreateParams struct {
	Runtime types.RuntimeProvider `json:"runtime"`
	SSHKey  types.SSHKey          `json:"sshKey"`
}

func (s *SessionCreateParams) Bind(r *http.Request) error {
	if s.Runtime == "" {
		return errors.New("field `runtime` is required")
	}
	if s.Runtime != types.LambdaLabsProvider && s.Runtime != types.UnweaveProvider {
		return fmt.Errorf("invalid runtime provider: %s. Must be one of `%s` or `%s`", s.Runtime, types.LambdaLabsProvider, types.UnweaveProvider)
	}
	return nil
}

func sessionsCreate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scr := SessionCreateParams{}
		if err := render.Bind(r, &scr); err != nil {
			log.Warn().
				Err(err).
				Msg("failed to read body")

			render.Render(w, r, ErrBadRequest("Invalid request body: "+err.Error()))
			return
		}

		rt, err := rti.FromUser(uuid.New(), scr.Runtime)
		node, err := rt.InitNode(r.Context(), scr.SSHKey)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("failed to init node")

			render.Render(w, r, ErrHTTPError(err, "Failed to initialize node"))
			return
		}

		// add to db
		res := &types.Session{ID: node.ID, SSHKey: node.KeyPair}
		render.JSON(w, r, res)
	}
}

func sessionsList(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := []*types.Session{
			{ID: "1"},
		}
		render.JSON(w, r, res)
	}
}

func sessionsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		res := &types.Session{ID: id}
		render.JSON(w, r, res)
	}
}
