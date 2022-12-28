package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/session"
	"github.com/unweave/unweave-v2/session/model"
)

type SessionCreateParams struct {
	Runtime model.RuntimeProvider `json:"runtime"`
	SSHKey  model.SSHKey          `json:"sshKey"`
}

func (s *SessionCreateParams) Bind(r *http.Request) error {
	if s.Runtime == "" {
		return errors.New("field `runtime` is required")
	}
	if s.Runtime != model.LambdaLabsProvider && s.Runtime != model.UnweaveProvider {
		return fmt.Errorf("invalid runtime provider: %s. Must be one of `%s` or `%s`", s.Runtime, model.LambdaLabsProvider, model.UnweaveProvider)
	}
	return nil
}

type Session struct {
	ID     string       `json:"id"`
	SSHKey model.SSHKey `json:"sshKey"`
	Status model.Status `json:"runtimeStatus"`
}

func sessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	scr := SessionCreateParams{}
	if err := render.Bind(r, &scr); err != nil {
		log.Warn().Err(err).Msg("failed to read body")
		render.Render(w, r, ErrBadRequest("Invalid request body: "+err.Error()))
		return
	}

	rt := session.NewRuntime(scr.Runtime)
	node, err := rt.InitNode(scr.SSHKey)
	if err != nil {
		log.Warn().Err(err).Msg("failed to init node")
		render.Render(w, r, ErrInternalServer("Failed to initialize node"))
		return
	}

	// add to db
	res := &Session{ID: node.ID, SSHKey: node.KeyPair}
	render.JSON(w, r, res)
}

func sessionGetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res := &Session{ID: id}
	render.JSON(w, r, res)
}
