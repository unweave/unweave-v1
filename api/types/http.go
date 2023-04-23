package types

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

const maxBuildContextSize = 1024 * 1024 * 100 // 100MB

func parseContextFile(r *http.Request) (multipart.File, error) {
	invalidFileErr := &Error{
		Code:       http.StatusBadRequest,
		Message:    "Failed to parse build context file",
		Suggestion: "Make sure the build context is a valid zipped multipart form file called 'context.zip'",
	}

	if err := r.ParseMultipartForm(maxBuildContextSize); err != nil {
		invalidFileErr.Err = fmt.Errorf("failed to parse multipart form: %w", err)
		return nil, invalidFileErr
	}

	// Only allowed to upload a single file called context.zip for now
	form := r.MultipartForm
	file, exists := form.File["context"]
	if !exists {
		invalidFileErr.Message = "No build context file found"
		return nil, invalidFileErr
	}
	if len(file) != 1 {
		invalidFileErr.Message = "More than one build context file found"
		return nil, invalidFileErr
	}
	if file[0].Filename != "context.zip" {
		invalidFileErr.Message = "Invalid build context file name"
		return nil, invalidFileErr
	}
	part, err := file[0].Open()
	if err != nil {
		invalidFileErr.Err = fmt.Errorf("failed to open file: %w", err)
		return nil, invalidFileErr
	}
	return part, nil
}

type BuildsCreateParams struct {
	Builder      string        `json:"builder"`
	Name         *string       `json:"name,omitempty"`
	BuildContext io.ReadCloser `json:"-"`
}

func (i *BuildsCreateParams) Bind(r *http.Request) error {
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
	part, err := parseContextFile(r)
	if err != nil {
		return err
	}
	i.BuildContext = part

	return nil
}

type BuildsCreateResponse struct {
	BuildID string `json:"buildID"`
}

type BuildsGetResponse struct {
	Build
	UsedBySessions []Exec      `json:"usedBySessions,omitempty"`
	Logs           *[]LogEntry `json:"logs,omitempty"`
}

type NodeTypesListResponse struct {
	NodeTypes []NodeType `json:"nodeTypes"`
}

type ExecCreateParams struct {
	Name          string   `json:"name,omitempty"`
	Provider      Provider `json:"provider"`
	NodeTypeID    string   `json:"nodeTypeID,omitempty"`
	Region        *string  `json:"region,omitempty"`
	SSHKeyName    *string  `json:"sshKeyName"`
	SSHPublicKey  *string  `json:"sshPublicKey"`
	IsInteractive bool     `json:"isInteractive"`
	PersistentFS  bool     `json:"persistentFS"`
	FilesystemID  *string  `json:"filesystemID,omitempty"`
	Ctx           ExecCtx  `json:"ctx,omitempty"`
}

func (s *ExecCreateParams) Bind(r *http.Request) error {
	// Check if header is set, if yes and set to json, parse body as json
	if r.Header.Get("Content-Type") == "application/json" {

		log.Info().Msg("Content-Type is application/json")

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(s); err != nil {
			return &Error{
				Code:       http.StatusBadRequest,
				Message:    "Failed to parse request body",
				Suggestion: "Make sure the request body is valid JSON",
				Err:        err,
			}
		}
		return nil
	}

	jsonStr := r.FormValue("params")
	if err := json.Unmarshal([]byte(jsonStr), s); err != nil {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    "Failed to parse request body",
			Suggestion: "Make sure the request body is a valid multipart form or valid JSON",
			Err:        err,
		}
	}
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

	if s.Ctx.Command != nil {
		part, err := parseContextFile(r)
		if err != nil {
			return err
		}
		s.Ctx.Context = part
	}

	return nil
}

type ProviderConnectParams struct {
	Provider      Provider `json:"provider"`
	ProviderToken string   `json:"providerToken,omitempty"`
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
	Providers []Provider `json:"providers"`
}

type SessionGetResponse struct {
	Session Exec `json:"session"`
}

type SessionsListResponse struct {
	Sessions []Exec `json:"sessions"`
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

type SSHKeyGenerateParams struct {
	Name *string `json:"name"`
}

func (s *SSHKeyGenerateParams) Bind(r *http.Request) error {
	return nil
}

type SSHKeyResponse struct {
	Name       string `json:"name"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SSHKeyListResponse struct {
	Keys []SSHKey `json:"keys"`
}
