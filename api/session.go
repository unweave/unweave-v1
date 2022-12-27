package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/unweave/unweave-v2/config"
)

// swagger:route POST /session/{id} session sessionCreate

// swagger:parameters sessionCreate
type sessionCreateRequest struct {
	// in: body
	Body SessionCreateRequest
}

type SessionCreateRequest struct {
	Runtime config.RuntimeProvider `json:"runtime"`
}

func (s *SessionCreateRequest) Bind(r *http.Request) error {
	if s.Runtime == "" {
		return errors.New("field `runtime` is required")
	}
	if s.Runtime != config.LambdaLabs && s.Runtime != config.Unweave {
		return fmt.Errorf("invalid runtime provider: %s. Must be one of `%s` or `%s`", s.Runtime, config.LambdaLabs, config.Unweave)
	}
	return nil
}

// swagger:response sessionCreateResponse
type sessionCreateResponse struct {
	// in: body
	Body SessionCreateResponse
}

type SessionCreateResponse struct {
	ID string `json:"id"`
}
