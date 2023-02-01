package types

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

type HTTPError struct {
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	Suggestion string          `json:"suggestion,omitempty"`
	Provider   RuntimeProvider `json:"provider,omitempty"`
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
		log.Ctx(r.Context()).Error().Err(e.Err).Stack().Msg(e.Message)
	} else {
		log.Ctx(r.Context()).Warn().Err(e.Err).Stack().Msg(e.Message)
	}
	render.Status(r, e.Code)
	return nil
}

type NodeTypesListResponse struct {
	NodeTypes []NodeType `json:"nodeTypes"`
}

type SessionCreateParams struct {
	Provider      RuntimeProvider `json:"provider"`
	NodeTypeID    string          `json:"nodeTypeID,omitempty"`
	ProviderToken *string         `json:"providerToken,omitempty"`
	Region        *string         `json:"region,omitempty"`
	SSHKeyName    *string         `json:"sshKeyName"`
	SSHPublicKey  *string         `json:"sshPublicKey"`
}

func (s *SessionCreateParams) Bind(r *http.Request) error {
	if s.Provider == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body: field 'provider' is required",
		}
	}
	if s.SSHPublicKey == nil && s.SSHKeyName == nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body: either 'sshKeyName' or 'sshPublicKey' is required",
		}
	}
	return nil
}

type ProviderConnectParams struct {
	Provider      RuntimeProvider `json:"provider"`
	ProviderToken string          `json:"providerToken,omitempty"`
}

func (p *ProviderConnectParams) Bind(r *http.Request) error {
	if p.Provider == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body: field 'provider' is required",
		}
	}
	return nil
}

type ProvidersListResponse struct {
	Providers []RuntimeProvider `json:"providers"`
}

type SessionGetResponse struct {
	Session Session `json:"session"`
}

type SessionsListResponse struct {
	Sessions []Session `json:"sessions"`
}

type SessionTerminateRequestParams struct {
	ProviderToken *string `json:"providerToken,omitempty"`
}

func (s *SessionTerminateRequestParams) Bind(r *http.Request) error {
	return nil
}

type SessionTerminateResponse struct {
	Success bool `json:"success"`
}

type SSHKeyAddParams struct {
	Name      *string `json:"name"`
	PublicKey string  `json:"publicKey"`
}

func (s *SSHKeyAddParams) Bind(r *http.Request) error {
	if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(s.PublicKey)); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid SSH public key",
		}
	}
	return nil
}

type SSHKeyAddResponse struct {
	Success bool `json:"success"`
}

type SSHKeyGenerateParams struct {
	Name *string `json:"name"`
}

func (s *SSHKeyGenerateParams) Bind(r *http.Request) error {
	return nil
}

type SSHKeyGenerateResponse struct {
	Name       string `json:"name"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SSHKeyListResponse struct {
	Keys []SSHKey `json:"keys"`
}
