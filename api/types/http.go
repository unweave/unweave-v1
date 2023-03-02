package types

import (
	"net/http"

	"golang.org/x/crypto/ssh"
)

const maxZeplContextSize = 1024 * 1024 * 100 // 100MB

type ImageBuildParams struct {
	ProjectID string `json:"projectID"`
}

func (i *ImageBuildParams) Bind(r *http.Request) error {
	return nil
}

type ImageBuildResponse struct {
	BuildID string `json:"buildID"`
}

type NodeTypesListResponse struct {
	NodeTypes []NodeType `json:"nodeTypes"`
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
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body: field 'provider' is required",
		}
	}
	if s.SSHPublicKey == nil && s.SSHKeyName == nil {
		return &Error{
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
		return &Error{
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

type SessionTerminateResponse struct {
	Success bool `json:"success"`
}

type SSHKeyAddParams struct {
	Name      *string `json:"name"`
	PublicKey string  `json:"publicKey"`
}

func (s *SSHKeyAddParams) Bind(r *http.Request) error {
	if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(s.PublicKey)); err != nil {
		return &Error{
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
