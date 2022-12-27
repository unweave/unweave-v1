package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave-v2/session"
	"github.com/unweave/unweave-v2/session/model"
)

// ---------------------------------------------------------------------------------------
// SessionCreate POST /session
// ---------------------------------------------------------------------------------------

// swagger:parameters sessionCreate
type sessionCreateRequest struct {
	// in: body
	Body SessionCreateRequest
}

type SessionCreateRequest struct {
	Runtime model.RuntimeProvider `json:"runtime"`
	model.SSHKey
}

func (s *SessionCreateRequest) Bind(r *http.Request) error {
	if s.Runtime == "" {
		return errors.New("field `runtime` is required")
	}
	if s.Runtime != model.LambdaLabsProvider && s.Runtime != model.UnweaveProvider {
		return fmt.Errorf("invalid runtime provider: %s. Must be one of `%s` or `%s`", s.Runtime, model.LambdaLabsProvider, model.UnweaveProvider)
	}
	return nil
}

// swagger:response sessionCreate
type sessionCreateResponse struct {
	// in: body
	Body SessionCreateResponse
}

type SessionCreateResponse struct {
	ID         string       `json:"id"`
	SSHKeyPair model.SSHKey `json:"sshKeyPair"`
}

func sessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	scr := SessionCreateRequest{}
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
	res := &SessionCreateResponse{ID: node.ID, SSHKeyPair: node.KeyPair}
	render.JSON(w, r, res)
}

// ---------------------------------------------------------------------------------------
// SessionGet GET /session/{id}
// ---------------------------------------------------------------------------------------

// swagger:response sessionGet
type sessionGetResponse struct {
	// in: body
	Body SessionGetResponse
}

type SessionGetResponse struct {
	ID     string       `json:"id"`
	Status model.Status `json:"runtimeStatus"`
}

// ---------------------------------------------------------------------------------------
// SessionConnect PUT /session/{id}/connect
// ---------------------------------------------------------------------------------------

// swagger:response sessionConnect
type sessionConnectResponse struct {
	// in: body
	Body SessionConnectResponse
}

type SessionConnectResponse struct {
	ID     string       `json:"id"`
	Status model.Status `json:"runtimeStatus"`
}
