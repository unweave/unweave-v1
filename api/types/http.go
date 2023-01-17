package types

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type HTTPError struct {
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	Suggestion string          `json:"suggestion"`
	Provider   RuntimeProvider `json:"provider"`
	Err        error           `json:"-"`
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *HTTPError) Render(w http.ResponseWriter, r *http.Request) error {
	// Depending on whether it is Unweave's fault or the user's fault, log the error
	// appropriately.
	if e.Code == http.StatusInternalServerError {
		log.Ctx(r.Context()).Error().Err(e).Stack().Msg(e.Message)
	} else {
		log.Ctx(r.Context()).Warn().Err(e).Stack().Msg(e.Message)
	}
	render.Status(r, e.Code)
	return nil
}

type SessionCreateParams struct {
	Provider     RuntimeProvider `json:"provider"`
	NodeTypeID   string          `json:"nodeTypeID,omitempty"`
	Region       *string         `json:"region,omitempty"`
	SSHKeyName   *string         `json:"sshKeyName"`
	SSHPublicKey *string         `json:"sshPublicKey"`
}

func (s *SessionCreateParams) Bind(r *http.Request) error {
	if s.Provider == "" {
		return &HTTPError{
			Code:       http.StatusBadRequest,
			Message:    "Invalid request body: field 'runtime' is required",
			Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", LambdaLabsProvider, UnweaveProvider),
		}
	}
	if s.Provider != LambdaLabsProvider && s.Provider != UnweaveProvider {
		return &HTTPError{
			Code:       http.StatusBadRequest,
			Message:    "Invalid runtime provider: " + string(s.Provider),
			Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", LambdaLabsProvider, UnweaveProvider),
		}
	}
	return nil
}

type SessionsListResponse struct {
	Sessions []Session `json:"sessions"`
}

type SessionTerminateResponse struct {
	Success bool `json:"success"`
}
