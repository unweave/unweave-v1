// Package api provides the API for the Unweave server.
//
// This file contains the types returned by the API.
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

// RuntimeProvider is the platform that the node is spawned on. This is where the user
// code runs
type RuntimeProvider string

func (r RuntimeProvider) String() string {
	return string(r)
}

type SessionStatus string

const (
	StatusInitializing SessionStatus = "initializing"
	StatusActive       SessionStatus = "active"
	StatusTerminated   SessionStatus = "terminated"
)

const (
	LambdaLabsProvider RuntimeProvider = "LambdaLabs"
	UnweaveProvider    RuntimeProvider = "Unweave"
)

type NodeSpecs struct {
	VCPUs int `json:"vCPUs"`
	// Memory is the RAM in GB
	Memory int `json:"memory"`
	// GPUMemory is the GPU RAM in GB
	GPUMemory *int `json:"gpuMemory"`
}

type NodeType struct {
	ID          string          `json:"id"`
	Name        *string         `json:"name"`
	Price       *int            `json:"price"`
	Regions     []string        `json:"regions"`
	Description *string         `json:"description"`
	Provider    RuntimeProvider `json:"provider"`
	Specs       NodeSpecs       `json:"specs"`
}

type Node struct {
	ID       string          `json:"id"`
	TypeID   string          `json:"typeID"`
	Region   string          `json:"region"`
	KeyPair  SSHKey          `json:"sshKeyPair"`
	Status   SessionStatus   `json:"status"`
	Provider RuntimeProvider `json:"provider"`
}

type SSHKey struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	PublicKey string    `json:"publicKey,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type ExecParams struct {
	Cmd   []string `json:"cmd"`
	Image []string `json:"containerImage"`
}

type Session struct {
	ID         uuid.UUID       `json:"id"`
	SSHKey     SSHKey          `json:"sshKey"`
	Status     SessionStatus   `json:"runtimeStatus"`
	NodeTypeID string          `json:"nodeTypeID"`
	Region     string          `json:"region"`
	Provider   RuntimeProvider `json:"provider"`
	NodeID     string          `json:"nodeID"`
}

type SessionCreateRequestParams struct {
	Provider     RuntimeProvider `json:"provider"`
	NodeTypeID   string          `json:"nodeTypeID,omitempty"`
	Region       *string         `json:"region,omitempty"`
	SSHKeyName   *string         `json:"sshKeyName"`
	SSHPublicKey *string         `json:"sshPublicKey"`
}

func (s *SessionCreateRequestParams) Bind(r *http.Request) error {
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

type SessionCreateResponse struct {
	Session Session `json:"session"`
}

type SessionsListResponse struct {
	Sessions []Session `json:"sessions"`
}

type SessionTerminateResponse struct {
	Success bool `json:"success"`
}

type SSHKeyAddRequestParams struct {
	Name      *string `json:"name"`
	PublicKey string  `json:"publicKey"`
}

func (s *SSHKeyAddRequestParams) Bind(r *http.Request) error {
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

type SSHKeyListResponse struct {
	Keys []SSHKey `json:"keys"`
}
