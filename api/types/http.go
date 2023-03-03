package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/ssh"
)

const maxBuildContextSize = 1024 * 1024 * 100 // 100MB

type ImageBuildParams struct {
	Builder      string        `json:"builder"`
	BuildContext io.ReadCloser `json:"-"`
}

func (i *ImageBuildParams) Bind(r *http.Request) error {
	jsonStr := r.FormValue("params")
	if err := json.Unmarshal([]byte(jsonStr), i); err != nil {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    "Failed to parse request body",
			Suggestion: "Make sure the request body is valid JSON",
			Err:        err,
		}
	}
	if i.Builder != "docker" {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    fmt.Sprintf("Invalid builder: %s", i.Builder),
			Suggestion: "Valid builders are: docker",
		}
	}

	// Validate build context in Multipart Form
	invalidFileErr := &Error{
		Code:       http.StatusBadRequest,
		Message:    "Failed to parse build context file",
		Suggestion: "Make sure the build context is a valid zipped multipart form file called 'context.zip'",
	}

	if err := r.ParseMultipartForm(maxBuildContextSize); err != nil {
		invalidFileErr.Err = fmt.Errorf("failed to parse multipart form: %w", err)
		return invalidFileErr
	}

	// Only allowed to upload a single file called context.zip for now
	form := r.MultipartForm
	file, exists := form.File["context"]
	if !exists {
		invalidFileErr.Message = "No build context file found"
		return invalidFileErr
	}
	if len(file) != 1 {
		invalidFileErr.Message = "More than one build context file found"
		return invalidFileErr
	}
	if file[0].Filename != "context.zip" {
		invalidFileErr.Message = "Invalid build context file name"
		return invalidFileErr
	}
	part, err := file[0].Open()
	if err != nil {
		invalidFileErr.Err = fmt.Errorf("failed to open file: %w", err)
		return invalidFileErr
	}
	i.BuildContext = part

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
